describe('Command execution', () => {
  it('should show retry attempts for a failing command', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Retry Test').click();

    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'Executing command: sh -c \'echo "attempt output"; exit 1\'', { timeout: 15000 }).should('exist');
    cy.contains('code', 'Retrying command (attempt 1/2)', { timeout: 15000 }).should('exist');
    cy.contains('code', 'Retrying command (attempt 2/2)', { timeout: 15000 }).should('exist');
    cy.contains('code', 'attempt output', { timeout: 15000 }).should('exist');
  });

  it('should cancel a hanging command on timeout', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Timeout Test').click();

    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'Executing command: sh -c \'echo "hanging"; sleep 30; echo "should not appear"\'', { timeout: 15000 }).should('exist');
    cy.contains('code', 'hanging', { timeout: 15000 }).should('exist');
    cy.contains('code', 'timed out', { timeout: 15000 }).should('exist');
    cy.contains('code', 'Job finished (took 2s)', { timeout: 15000 }).should('exist');
  });

  it('should stop on first failing command with fail_fast', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Fail Fast Stops').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'Executing command: echo "before fail"', { timeout: 15000 }).should('exist');
    cy.contains('code', 'before fail', { timeout: 15000 }).should('exist');
    cy.contains('code', 'Executing command: sh -c \'echo "failing"; exit 1\'', { timeout: 15000 }).should('exist');
    cy.contains('code', 'failing', { timeout: 15000 }).should('exist');
    cy.contains('code', 'exit status 1', { timeout: 15000 }).should('exist');
    cy.contains('code', 'Executing command: echo "after fail should not appear"', { timeout: 15000 }).should('not.exist');
  });

  it('should continue all commands with disable_fail_fast', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Continue On Failure').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'Executing command: echo "before fail"', { timeout: 15000 }).should('exist');
    cy.contains('code', 'before fail', { timeout: 15000 }).should('exist');
    cy.contains('code', 'Executing command: sh -c \'echo "failing"; exit 1\'', { timeout: 15000 }).should('exist');
    cy.contains('code', 'failing', { timeout: 15000 }).should('exist');
    cy.contains('code', 'exit status 1', { timeout: 15000 }).should('exist');
    cy.contains('code', 'Executing command: echo "after fail continues"', { timeout: 15000 }).should('exist');
    cy.contains('code', 'after fail continues', { timeout: 15000 }).should('exist');
  });

  it('should execute pre and post commands from job_defaults', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Pre And Post Commands').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'Executing command: echo "Starting backup..."', { timeout: 15000 }).should('exist');
    cy.contains('code', 'Starting backup...', { timeout: 15000 }).should('exist');
    cy.contains('code', 'Executing command: echo "job command"', { timeout: 15000 }).should('exist');
    cy.contains('code', 'job command', { timeout: 15000 }).should('exist');
    cy.contains('code', 'Executing command: echo "Backup finished!"', { timeout: 15000 }).should('exist');
    cy.contains('code', 'Backup finished!', { timeout: 15000 }).should('exist');
  });

  it('should retry a timed out command until it succeeds', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Timeout Then Retry Succeeds').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'Retrying command (attempt 1/2)', { timeout: 15000 }).should('exist');
    cy.contains('code', 'attempt 2 succeeded', { timeout: 15000 }).should('exist');
    cy.contains('code', 'Job finished', { timeout: 15000 }).should('exist');
  });

  it('should exhaust retries when every attempt times out', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E All Retries Timeout').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'always times out', { timeout: 15000 }).should('exist');
    cy.contains('code', 'timed out', { timeout: 15000 }).should('exist');
    cy.contains('code', 'Retrying command (attempt 1/1)', { timeout: 15000 }).should('exist');
    cy.contains('code', 'Job finished', { timeout: 15000 }).should('exist');
  });

  it('should capture stderr output', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Stderr Captured').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'to stdout', { timeout: 15000 }).should('exist');
    cy.contains('code', 'to stderr', { timeout: 15000 }).should('exist');
  });

  it('should show no output message for silent commands', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Empty Output').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'No output', { timeout: 15000 }).should('exist');
  });
});
