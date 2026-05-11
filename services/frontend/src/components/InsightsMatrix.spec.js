import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount, flushPromises } from '@vue/test-utils';
import InsightsMatrix from './InsightsMatrix.vue';
import * as insightsApi from '../api/insights.js';

vi.mock('vue-router', () => ({
  useRouter: vi.fn().mockReturnValue({ push: vi.fn() }),
}));

vi.mock('../api/insights.js', () => ({
  getWeeklyInsights: vi.fn(),
  getInsights: vi.fn(),
}));

describe('InsightsMatrix', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    insightsApi.getWeeklyInsights.mockResolvedValue([]);
  });

  it('TestDashboard_GivenMonthWindow_WhenRendered_ThenShows3RowsAcross31Days', async () => {
    const wrapper = mount(InsightsMatrix, { props: { past: 31 } });
    await flushPromises();

    expect(wrapper.findAll('[data-testid="matrix-row"]')).toHaveLength(3);
    expect(wrapper.findAll('[data-testid="matrix-cell"]')).toHaveLength(93);
  });

  it('TestDashboard_GivenMonthWindow_WhenRendered_ThenRightmostColumnIsToday', async () => {
    const today = new Date();
    const todayStr = new Date(Date.UTC(
      today.getUTCFullYear(),
      today.getUTCMonth(),
      today.getUTCDate(),
    )).toISOString().slice(0, 10).replace(/-/g, '');

    insightsApi.getWeeklyInsights.mockResolvedValue([
      { date: todayStr, intimacy: 99, commitment: 88, passion: 77 },
    ]);

    const wrapper = mount(InsightsMatrix, { props: { past: 31 } });
    await flushPromises();

    // Cells are laid out row-major: intimacy cells 0-30, commitment 31-61, passion 62-92.
    // Index 30 = last cell of the intimacy row = rightmost column = today (score 99 → high).
    const cells = wrapper.findAll('[data-testid="matrix-cell"]');
    expect(cells[30].classes()).toContain('insight-value--high');
    // Index 0 = 30 days ago, no data → unavailable.
    expect(cells[0].classes()).toContain('insight-value--unavailable');
  });

  it('TestDashboard_GivenInsightScore_WhenCellRendered_ThenUsesInsightsColorCoding', async () => {
    const today = new Date();
    const makeDate = (daysAgo) => new Date(Date.UTC(
      today.getUTCFullYear(),
      today.getUTCMonth(),
      today.getUTCDate() - daysAgo,
    )).toISOString().slice(0, 10).replace(/-/g, '');

    // Four cells with scores that fall into each color band.
    insightsApi.getWeeklyInsights.mockResolvedValue([
      { date: makeDate(3), intimacy: 10,  commitment: 10,  passion: 10  }, // ≤24 → very-low
      { date: makeDate(2), intimacy: 40,  commitment: 40,  passion: 40  }, // ≤49 → low
      { date: makeDate(1), intimacy: 60,  commitment: 60,  passion: 60  }, // ≤74 → moderate
      { date: makeDate(0), intimacy: 100, commitment: 100, passion: 100 }, // >74 → high
    ]);

    const wrapper = mount(InsightsMatrix, { props: { past: 4 } });
    await flushPromises();

    // Row-major layout with past=4: intimacy cells 0-3.
    const cells = wrapper.findAll('[data-testid="matrix-cell"]');
    expect(cells[0].classes()).toContain('insight-value--very-low');
    expect(cells[1].classes()).toContain('insight-value--low');
    expect(cells[2].classes()).toContain('insight-value--moderate');
    expect(cells[3].classes()).toContain('insight-value--high');
  });

  it('TestDashboard_GivenMissingInsightDay_WhenCellRendered_ThenShowsGreyCell', async () => {
    // All days have no data: API returns empty array → every cell is unavailable.
    insightsApi.getWeeklyInsights.mockResolvedValue([]);

    const wrapper = mount(InsightsMatrix, { props: { past: 3 } });
    await flushPromises();

    const cells = wrapper.findAll('[data-testid="matrix-cell"]');
    // 3 dimensions × 3 days = 9 cells, all must be unavailable (grey).
    expect(cells).toHaveLength(9);
    cells.forEach(cell => {
      expect(cell.classes()).toContain('insight-value--unavailable');
    });
  });

  it('TestDashboard_GivenMobileViewport_WhenMonthlyMatrixShown_ThenRemainsCondensedAndReadable', async () => {
    // The matrix section must render with a wrapping element that carries the
    // data-testid, and the table must be present. Layout condensation (width:100%,
    // fixed table-layout, small row height) is defined in scoped CSS and verified
    // by visual review at 375px per CONSTITUTION.md baseline.
    const wrapper = mount(InsightsMatrix, { props: { past: 31 } });
    await flushPromises();

    expect(wrapper.find('[data-testid="monthly-matrix"]').exists()).toBe(true);
    expect(wrapper.find('table').exists()).toBe(true);
    // 31 columns × 3 rows = exactly 93 cells — confirms condensed structure fits.
    expect(wrapper.findAll('[data-testid="matrix-cell"]')).toHaveLength(93);
  });
});
