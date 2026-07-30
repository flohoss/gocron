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
});
