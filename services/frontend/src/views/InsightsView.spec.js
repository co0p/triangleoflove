import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount, flushPromises } from '@vue/test-utils';
import InsightsView from './InsightsView.vue';
import * as insightsApi from '../api/insights.js';

const { push } = vi.hoisted(() => ({
  push: vi.fn(),
}));

vi.mock('vue-router', () => ({
  useRoute: vi.fn().mockReturnValue({ params: { date: '20260330' } }),
  useRouter: vi.fn().mockReturnValue({ push }),
}));

vi.mock('../api/insights.js', () => ({
  getInsights: vi.fn(),
}));

describe('InsightsView', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    insightsApi.getInsights.mockResolvedValue({
      intimacy: 100,
      commitment: 100,
      passion: 100,
    });
  });

  it('TestInsightsView_GivenScoresLoaded_WhenRendered_ThenRawValuesVisible', async () => {
    const wrapper = mount(InsightsView);
    await flushPromises();

    expect(insightsApi.getInsights).toHaveBeenCalledWith('20260330');
    expect(wrapper.get('[data-testid="insight-intimacy"]').text()).toBe('100');
    expect(wrapper.get('[data-testid="insight-commitment"]').text()).toBe('100');
    expect(wrapper.get('[data-testid="insight-passion"]').text()).toBe('100');
  });

  it('TestInsightsView_GivenScoresLoaded_WhenRendered_ThenColorsMatchBands', async () => {
    const scenarios = [
      {
        scores: { intimacy: 24, commitment: 24, passion: 24 },
        expectedClass: 'insight-value--very-low',
      },
      {
        scores: { intimacy: 49, commitment: 49, passion: 49 },
        expectedClass: 'insight-value--low',
      },
      {
        scores: { intimacy: 74, commitment: 74, passion: 74 },
        expectedClass: 'insight-value--moderate',
      },
      {
        scores: { intimacy: 100, commitment: 100, passion: 100 },
        expectedClass: 'insight-value--high',
      },
      {
        scores: { intimacy: -1, commitment: -1, passion: -1 },
        expectedClass: 'insight-value--unavailable',
      },
    ];

    for (const scenario of scenarios) {
      insightsApi.getInsights.mockResolvedValueOnce(scenario.scores);

      const wrapper = mount(InsightsView);
      await flushPromises();

      expect(wrapper.get('[data-testid="insight-intimacy"]').classes()).toContain(scenario.expectedClass);
      expect(wrapper.get('[data-testid="insight-commitment"]').classes()).toContain(scenario.expectedClass);
      expect(wrapper.get('[data-testid="insight-passion"]').classes()).toContain(scenario.expectedClass);

      wrapper.unmount();
    }
  });

  it('TestInsightsView_GivenNoData_WhenRendered_ThenShowsNoDataMessage', async () => {
    insightsApi.getInsights.mockRejectedValue(new Error('not_found'));

    const wrapper = mount(InsightsView);
    await flushPromises();

    expect(wrapper.get('[role="alert"]').text()).toContain('No insights available for this date.');
  });

  it('TestInsightsView_GivenInvalidDateError_WhenRendered_ThenShowsErrorMessage', async () => {
    insightsApi.getInsights.mockRejectedValue(new Error('invalid_date'));

    const wrapper = mount(InsightsView);
    await flushPromises();

    expect(wrapper.get('[role="alert"]').text()).toContain('Invalid date. Use YYYYMMDD.');
  });
});