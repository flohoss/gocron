describe('should render the GUI', () => {
  it('should render the jobs page', () => {
    cy.visit('/').then(() => {
      cy.get('.container').should('be.visible');

      // Left side actions
      cy.get('[data-test-id="terminal-button"]').should('be.visible');
      cy.get('[data-test-id="openapi-button"]').should('be.visible');

      // Right side actions
      cy.get('[data-test-id="select-button"]').should('be.visible');
      cy.get('[data-test-id="run-button"]').should('be.visible').should('have.attr', 'data-tip', 'Run All Jobs');

      // Detail page action
      cy.get('[data-test-id="back-button"]').should('not.exist');
    });
  });

  it('should render the job detail page', () => {
    cy.visit('/');

    cy.get('[data-test-id="job-link"]').should('be.visible');

    cy.get('[data-test-id="job-link"]')
      .its('length')
      .then((jobCount) => {
        Cypress._.times(jobCount, (index) => {
          cy.get('[data-test-id="job-link"]')
            .eq(index)
            .then(($link) => {
              const jobName = $link.attr('data-test-name');

              cy.wrap($link).click({ force: true });

              // Left side action
              cy.get('[data-test-id="back-button"]').should('be.visible');

              // Right side action
              cy.get('[data-test-id="run-button"]').should('be.visible').should('have.attr', 'data-tip', `Run ${jobName}`);

              // Home page actions
              cy.get('[data-test-id="terminal-button"]').should('not.exist');
              cy.get('[data-test-id="openapi-button"]').should('not.exist');
              cy.get('[data-test-id="select-button"]').should('not.exist');

              cy.go('back');
              cy.location('pathname').should('eq', '/');
              cy.get('[data-test-id="job-link"]').should('be.visible');
            });
        });
      });
  });

  it('should navigate to command view and back', () => {
    cy.visit('/');

    cy.get('[data-test-id="terminal-button"]').click();
    cy.url().should('include', '/commands');

    cy.get('[data-test-id="back-button"]').should('be.visible').click();
    cy.url().should('match', /\/$/);
    cy.get('[data-test-id="terminal-button"]').should('be.visible');
  });

  it('should expose the openapi docs link', () => {
    cy.visit('/');

    cy.get('[data-test-id="openapi-button"] a').should('have.attr', 'href', '/api/docs');
  });

  it('should open and close the select jobs modal', () => {
    cy.visit('/');

    cy.get('[data-test-id="select-button"] button').click();
    cy.get('#selectModal').should('have.attr', 'open');
    cy.contains('h3', 'Select Jobs').should('be.visible');
    cy.get('input[type="search"]').should('be.visible').type('job');

    cy.get('form.modal-backdrop button').click({ force: true });
    cy.get('#selectModal').should('not.have.attr', 'open');
  });

  it('should reject a disallowed terminal command', () => {
    cy.visit('/commands');

    cy.get('input[placeholder="Command"]').type('whoami{enter}');
    cy.contains('code', 'Executing command: whoami').should('be.visible');
    cy.contains('code', 'command "whoami" is not allowed').should('be.visible');
  });

  it('should reject disallowed arguments for allowed command', () => {
    cy.visit('/commands');

    cy.get('input[placeholder="Command"]').type('docker images{enter}');
    cy.contains('code', 'Executing command: docker images').should('be.visible');
    cy.contains('code', 'argument "images" is not allowed for command "docker"').should('be.visible');
  });

  it('should show retry attempts for a failing command', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Retry Test').click();

    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'Executing command: sh -c \'echo "attempt output"; exit 1\'', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'Retrying command (attempt 1/2)', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'Retrying command (attempt 2/2)', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'attempt output', { timeout: 15000 }).should('be.visible');
  });

  it('should cancel a hanging command on timeout', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Timeout Test').click();

    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'Executing command: sh -c \'echo "hanging"; sleep 30; echo "should not appear"\'', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'hanging', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'timed out', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'Job finished (took 2s)', { timeout: 15000 }).should('be.visible');
  });

  it('should inherit default env vars from job_defaults', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Default Env Inherited').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'Setting environment variables: SLEEP_TIME', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'sleep=5', { timeout: 15000 }).should('be.visible');
  });

  it('should override default env with job-specific env', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Job Env Overrides Default').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'Setting environment variables: SLEEP_TIME', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'sleep=99', { timeout: 15000 }).should('be.visible');
  });

  it('should execute pre and post commands from job_defaults', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Pre And Post Commands').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'Executing command: echo "Starting backup..."', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'Starting backup...', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'Executing command: echo "job command"', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'job command', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'Executing command: echo "Backup finished!"', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'Backup finished!', { timeout: 15000 }).should('be.visible');
  });

  it('should stop on first failing command with fail_fast', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Fail Fast Stops').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'Executing command: echo "before fail"', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'before fail', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'Executing command: sh -c \'echo "failing"; exit 1\'', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'failing', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'exit status 1', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'Executing command: echo "after fail should not appear"', { timeout: 15000 }).should('not.exist');
  });

  it('should continue all commands with disable_fail_fast', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Continue On Failure').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'Executing command: echo "before fail"', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'before fail', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'Executing command: sh -c \'echo "failing"; exit 1\'', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'failing', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'exit status 1', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'Executing command: echo "after fail continues"', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'after fail continues', { timeout: 15000 }).should('be.visible');
  });

  it('should expand multiple env vars in commands', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Multiple Envs').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'Setting environment variables: SLEEP_TIME, BACKUP_DIR, RETENTION', { timeout: 15000 }).should('be.visible');
    cy.contains('code', 'dir=/tmp/backup retention=7', { timeout: 15000 }).should('be.visible');
  });
});
