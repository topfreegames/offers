Feature: Healthcheck

  Scenario Outline: Healthcheck verifies services
    Given the server is "<state>"
    When the health check is done
    Then the last request returned status code <status> and body "<body>"

    Examples:
      | state       | status | body                                                                    |
      | healthy     | 200    | { 'success': true, 'services': { 'postgres': true, 'redis': true } }    |
      | no-postgres | 500    | { 'success': false, 'services': { 'postgres': false, 'redis': true } }  |
      | no-redis    | 500    | { 'success': false, 'services': { 'postgres': true, 'redis': false } }  |
      | chaos       | 500    | { 'success': false, 'services': { 'postgres': false, 'redis': false } } |
