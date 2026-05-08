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
});