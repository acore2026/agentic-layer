Feature: End-to-End System Verification
  In order to ensure intent-driven networking functions correctly
  As a 6G network operator
  I need the system to accurately translate natural language to network operations

  Scenario: Successful Fleet Wake-Up
    Given all agentic core services are running
    When I send the intent "Wake up the fleet for firmware updates"
    Then the system should trigger the mcp://skill/device/fleet-update skill and return success

  Scenario: Unregistered Skill Discovery
    Given all agentic core services are running
    When I send the intent "Cook a pizza"
    Then the system should fail gracefully with a not found message
