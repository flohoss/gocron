describe('Server-sent events', () => {
  it('should deliver the run output over the event stream', () => {
    // Guards the live log view: the job page is only updated by SSE, so the
    // final command output has to arrive on the stream, not just in the API.
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Write And Read File').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'backup content', { timeout: 15000 }).should('exist');
    cy.contains('code', 'Job finished', { timeout: 15000 }).should('exist');
  });

  it('should emit events while the job is still running', () => {
    // Each log line is published on its own event, so early output must appear
    // before the run has finished rather than all at once at the end.
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Retry Test').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'Executing command: echo "Starting backup..."', { timeout: 15000 }).should('exist');
    cy.contains('code', 'Retrying command (attempt 1/2)', { timeout: 15000 }).should('exist');
  });

  it('should keep streaming across several runs on one connection', () => {
    // Repeated publish cycles on the same stream must keep working: catches a
    // stream that stops delivering or leaks subscribers after a few runs.
    const jobs = [
      ['E2E Default Env Inherited', 'sleep=5'],
      ['E2E Job Env Overrides Default', 'sleep=99'],
      ['E2E Multiple Envs', 'dir=/tmp/backup retention=7'],
    ];

    cy.visit('/');

    jobs.forEach(([jobName, expected]) => {
      cy.contains('[data-test-id="job-link"]', jobName).click();
      cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();
      cy.contains('code', expected, { timeout: 15000 }).should('exist');
      cy.contains('code', 'Job finished', { timeout: 15000 }).should('exist');
      cy.get('[data-test-id="back-button"]').click();
    });
  });

  it('should stream to every subscriber at once', () => {
    // Regression test for the broadcasting bug: a second subscriber used to be
    // able to stall delivery for all others. An extra idle EventSource must not
    // stop the page from receiving updates.
    cy.visit('/', {
      onBeforeLoad(win) {
        win.__streamUrl = null;
        const NativeEventSource = win.EventSource;
        win.EventSource = function (...args) {
          win.__streamUrl = win.__streamUrl || args[0];
          return new NativeEventSource(...args);
        };
      },
    });

    cy.contains('[data-test-id="job-link"]', 'E2E Pipe Between Commands').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();
    cy.contains('code', 'hello world', { timeout: 15000 }).should('exist');

    cy.window().then((win) => {
      win.__extraEvents = 0;
      win.__extraSource = new win.EventSource(win.__streamUrl);
      win.__extraSource.addEventListener('message', () => {
        win.__extraEvents++;
      });
    });

    cy.get('[data-test-id="back-button"]').click();
    cy.contains('[data-test-id="job-link"]', 'E2E Delete File').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'deleted', { timeout: 15000 }).should('exist');
    cy.window().its('__extraEvents').should('be.greaterThan', 0);

    cy.window().then((win) => {
      win.__extraSource.close();
    });
  });
});
