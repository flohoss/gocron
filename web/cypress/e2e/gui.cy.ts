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

    cy.get('[data-test-id="job-link"]').each(($link) => {
      const jobName = $link.attr('data-test-name');

      cy.wrap($link).click();

      // Left side action
      cy.get('[data-test-id="back-button"]').should('be.visible');

      // Right side action
      cy.get('[data-test-id="run-button"]').should('be.visible').should('have.attr', 'data-tip', `Run ${jobName}`);

      // Home page actions
      cy.get('[data-test-id="terminal-button"]').should('not.exist');
      cy.get('[data-test-id="openapi-button"]').should('not.exist');
      cy.get('[data-test-id="select-button"]').should('not.exist');

      cy.go('back');
      cy.get('[data-test-id="job-link"]').should('be.visible');
    });
  });
});
