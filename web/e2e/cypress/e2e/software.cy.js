describe('Software install', () => {
  it('should install sqlite3 at startup and report its version', () => {
    cy.visit('/commands');

    cy.get('input[placeholder="Command"]').type('sqlite3 --version{enter}');
    cy.contains('code', 'Executing command: sqlite3 --version').should('be.visible');
    cy.contains('code', '3.', { timeout: 10000 }).should('be.visible');
  });

  it('should skip already installed software and leave it usable', () => {
    cy.visit('/commands');

    cy.get('input[placeholder="Command"]').type('git --version{enter}');
    cy.contains('code', 'Executing command: git --version').should('be.visible');
    cy.contains('code', 'git version', { timeout: 10000 }).should('be.visible');
  });

  it('should install the exact pinned apprise version', () => {
    cy.visit('/commands');

    cy.get('input[placeholder="Command"]').type('apprise --version{enter}');
    cy.contains('code', 'Executing command: apprise --version').should('be.visible');
    cy.contains('code', 'Apprise v1.9.0', { timeout: 10000 }).should('be.visible');
  });
});
