describe('Shell behavior', () => {
  it('should write and read a file', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Write And Read File').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'backup content', { timeout: 15000 }).should('exist');
  });

  it('should create nested directory structures', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Create Directory Structure').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'created', { timeout: 15000 }).should('exist');
  });

  it('should delete a file and confirm removal', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Delete File').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'deleted', { timeout: 15000 }).should('exist');
  });

  it('should pipe output between commands', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Pipe Between Commands').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'hello world', { timeout: 15000 }).should('exist');
  });

  it('should execute command substitution', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Command Substitution').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'word count: 3', { timeout: 15000 }).should('exist');
  });

  it('should make a script executable and run it', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Make Executable And Run').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'script ran', { timeout: 15000 }).should('exist');
  });

  it('should continue after a specific exit code with disable_fail_fast', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Specific Exit Code').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'continued after exit 42', { timeout: 15000 }).should('exist');
  });

  it('should expand env vars inside command substitution', () => {
    cy.visit('/');

    cy.contains('[data-test-id="job-link"]', 'E2E Env In Command Substitution').click();
    cy.get('[data-test-id="run-button"]').should('not.be.disabled').click();

    cy.contains('code', 'file=nightly-2026', { timeout: 15000 }).should('exist');
  });
});
