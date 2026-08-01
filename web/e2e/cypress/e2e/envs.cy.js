describe('Environment variables', () => {
  it('should inherit default env vars from job_defaults', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Default Env Inherited').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'Setting environment variables: SLEEP_TIME', { timeout: 15000 }).should('exist');
    cy.contains('code', 'sleep=5', { timeout: 15000 }).should('exist');
  });

  it('should override default env with job-specific env', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Job Env Overrides Default').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'Setting environment variables: SLEEP_TIME', { timeout: 15000 }).should('exist');
    cy.contains('code', 'sleep=99', { timeout: 15000 }).should('exist');
  });

  it('should expand multiple env vars in commands', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Multiple Envs').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'Setting environment variables: SLEEP_TIME, BACKUP_DIR, RETENTION', { timeout: 15000 }).should('exist');
    cy.contains('code', 'dir=/tmp/backup retention=7', { timeout: 15000 }).should('exist');
  });

  it('should expand env vars in commands with job-specific env', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Env In Pre Commands').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'Setting environment variables: SLEEP_TIME, GREETING', { timeout: 15000 }).should('exist');
    cy.contains('code', 'hello from env', { timeout: 15000 }).should('exist');
  });

  it('should expand env vars without braces using $VAR syntax', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Env No Braces Syntax').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'Setting environment variables: SLEEP_TIME, SIMPLE_VAR', { timeout: 15000 }).should('exist');
    cy.contains('code', 'value=no-braces', { timeout: 15000 }).should('exist');
  });

  it('should handle env values with special characters', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Env Special Characters').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'Setting environment variables: SLEEP_TIME, SPECIAL', { timeout: 15000 }).should('exist');
    cy.contains('code', 'special=[has-spaces-and-dashes]', { timeout: 15000 }).should('exist');
  });
});
