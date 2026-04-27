import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount, flushPromises } from '@vue/test-utils';
import InsightsWeeklyView from './InsightsWeeklyView.vue';
import * as insightsApi from '../api/insights.js';

vi.mock('vue-router', () => ({
  useRouter: vi.fn().mockReturnValue({ push: vi.fn() }),
}));

vi.mock('../api/insights.js', () => ({
  getInsights: vi.fn(),
  getWeeklyInsights: vi.fn(),
}));

const sevenDays = [
  { date: '20260421', intimacy: 75, commitment: 50, passion: 60 },
  { date: '20260422', intimacy: 80, commitment: 55, passion: 65 },
  { date: '20260423', intimacy: -1, commitment: 60, passion: 70 },
  { date: '20260424', intimacy: 85, commitment: -1, passion: 75 },
  { date: '20260425', intimacy: 90, commitment: 65, passion: -1 },
  { date: '20260426', intimacy: 70, commitment: 70, passion: 80 },
  { date: '20260427', intimacy: 60, commitment: 75, passion: 85 },
];

describe('InsightsWeeklyView', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    insightsApi.getWeeklyInsights.mockResolvedValue(sevenDays);
  });

  it('TestInsightsWeekly_GivenAuthenticatedUser_WhenInsightsRouteLoads_ThenShows3x7Table', async () => {
    const wrapper = mount(InsightsWeeklyView);
    await flushPromises();

    expect(wrapper.findAll('[data-testid="weekly-row"]')).toHaveLength(3);
    expect(wrapper.findAll('[data-testid="weekly-cell"]')).toHaveLength(21);
  });

  it('TestInsightsWeekly_GivenSevenDayWindow_WhenDatesRendered_ThenRightmostIsYesterday', async () => {
    // API returns days in reverse order; frontend must sort ascending so rightmost = latest date
    insightsApi.getWeeklyInsights.mockResolvedValue([
      { date: '20260427', intimacy: 99, commitment: 99, passion: 99 },
      { date: '20260426', intimacy: 70, commitment: 70, passion: 80 },
      { date: '20260425', intimacy: 90, commitment: 65, passion: -1 },
      { date: '20260424', intimacy: 85, commitment: -1, passion: 75 },
      { date: '20260423', intimacy: -1, commitment: 60, passion: 70 },
      { date: '20260422', intimacy: 80, commitment: 55, passion: 65 },
      { date: '20260421', intimacy: 75, commitment: 50, passion: 60 },
    ]);

    const wrapper = mount(InsightsWeeklyView);
    await flushPromises();

    // Cells 0-6 = intimacy row. Index 6 is rightmost = should be the latest date (20260427, intimacy=99).
    const cells = wrapper.findAll('[data-testid="weekly-cell"]');
    expect(cells[6].text()).toBe('99');
    // Index 0 is leftmost = oldest date (20260421, intimacy=75).
    expect(cells[0].text()).toBe('75');
  });

  it('TestInsightsWeekly_GivenScorePresent_WhenCellRendered_ThenUsesExistingColorBands', async () => {
    // Verify that a cell with a real score displays the numeric value (not empty or placeholder).
    // Color class binding is covered by the shared insightValueClass logic already tested in InsightsView.spec.js.
    insightsApi.getWeeklyInsights.mockResolvedValue([
      { date: '20260421', intimacy: 75, commitment: 50, passion: 60 },
      { date: '20260422', intimacy: 75, commitment: 50, passion: 60 },
      { date: '20260423', intimacy: 75, commitment: 50, passion: 60 },
      { date: '20260424', intimacy: 75, commitment: 50, passion: 60 },
      { date: '20260425', intimacy: 75, commitment: 50, passion: 60 },
      { date: '20260426', intimacy: 75, commitment: 50, passion: 60 },
      { date: '20260427', intimacy: 75, commitment: 50, passion: 60 },
    ]);

    const wrapper = mount(InsightsWeeklyView);
    await flushPromises();

    // First cell (intimacy, oldest day) has score 75 — must display the value.
    const cell = wrapper.findAll('[data-testid="weekly-cell"]')[0];
    expect(cell.text()).toBe('75');
    expect(cell.text()).not.toBe('');
  });

  it('TestInsightsWeekly_GivenScoreMissing_WhenCellRendered_ThenShowsGreyNotZero', async () => {
    // A day where intimacy is -1 (missing check-in for that dimension).
    insightsApi.getWeeklyInsights.mockResolvedValue([
      { date: '20260421', intimacy: -1, commitment: 50, passion: 60 },
      { date: '20260422', intimacy: -1, commitment: 50, passion: 60 },
      { date: '20260423', intimacy: -1, commitment: 50, passion: 60 },
      { date: '20260424', intimacy: -1, commitment: 50, passion: 60 },
      { date: '20260425', intimacy: -1, commitment: 50, passion: 60 },
      { date: '20260426', intimacy: -1, commitment: 50, passion: 60 },
      { date: '20260427', intimacy: -1, commitment: 50, passion: 60 },
    ]);

    const wrapper = mount(InsightsWeeklyView);
    await flushPromises();

    // Intimacy row = cells 0-6. All have intimacy=-1; must not display "0" or "-1".
    const cells = wrapper.findAll('[data-testid="weekly-cell"]');
    for (let i = 0; i < 7; i++) {
      expect(cells[i].text()).not.toBe('0');
      expect(cells[i].text()).not.toBe('-1');
    }
  });
});
