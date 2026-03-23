import assert from 'node:assert/strict';
import test from 'node:test';
import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';

test('TestRailwayDeploy_GivenMainPush_WhenCiAndApiPass_ThenDeployTriggered', () => {
  const orchestratorPath = resolve(process.cwd(), '..', '.github', 'workflows', 'ci-orchestrator.yml');
  const workflow = readFileSync(orchestratorPath, 'utf8');

  assert.match(workflow, /\n  deployment:\n/);
  assert.match(workflow, /\n    needs:\n      - api-tests\n/);
});

test('TestRailwayDeploy_GivenNonMainEvent_WhenCiRuns_ThenNoProductionDeploy', () => {
  const orchestratorPath = resolve(process.cwd(), '..', '.github', 'workflows', 'ci-orchestrator.yml');
  const workflow = readFileSync(orchestratorPath, 'utf8');

  assert.match(workflow, /\n  deployment:\n/);
  assert.match(workflow, /\n    if: github\.event_name == 'push' && github\.ref == 'refs\/heads\/main'\n/);
});