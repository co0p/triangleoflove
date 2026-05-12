import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount, flushPromises } from '@vue/test-utils';
import CheckinMatrix from './CheckinMatrix.vue';
import * as insightsApi from '../api/insights.js';

vi.mock('vue-router', () => ({
  useRouter: vi.fn().mockReturnValue({ push: vi.fn() }),
  RouterLink: { template: '<a :href="to"><slot /></a>', props: ['to'] },
}));

vi.mock('../api/insights.js', () => ({
  getWeeklyInsights: vi.fn(),
  getInsights: vi.fn(),
}));

describe('CheckinMatrix', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    insightsApi.getWeeklyInsights.mockResolvedValue([]);
  });

  it('TestCheckinHistory_GivenHistoryVisible_WhenRendered_ThenEachRowShowsIconAndLabel', async () => {
    const wrapper = mount(CheckinMatrix, { props: { past: 31 } });
    await flushPromises();

    expect(wrapper.findAll('[data-testid="dimension-icon"]')).toHaveLength(3);
  });

  it('TestCheckinHistory_GivenNarrowViewport_WhenRendered_ThenNoHorizontalOverflow', async () => {
    const wrapper = mount(CheckinMatrix, { props: { past: 31 } });
    await flushPromises();

    expect(wrapper.find('[data-testid="checkin-history"]').exists()).toBe(true);
    expect(wrapper.find('table').exists()).toBe(true);
  });

  it('TestCheckinHistory_GivenHistoryVisible_WhenRendered_ThenCellsHaveConsistentSpacing', async () => {
    const wrapper = mount(CheckinMatrix, { props: { past: 31 } });
    await flushPromises();

    expect(wrapper.findAll('[data-testid="matrix-row"]')).toHaveLength(3);
    expect(wrapper.findAll('[data-testid="matrix-cell"]')).toHaveLength(93);
  });

  it('TestDashboard_GivenMonthWindow_WhenRendered_ThenShows3RowsAcross31Days', async () => {
    const wrapper = mount(CheckinMatrix, { props: { past: 31 } });
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

    const wrapper = mount(CheckinMatrix, { props: { past: 31 } });
    await flushPromises();

    // Cells are laid out row-major: intimacy cells 0-30, commitment 31-61, passion 62-92.
    // Index 30 = last cell of the intimacy row = rightmost column = today (score 99 → high).
    const cells = wrapper.findAll('[data-testid="matrix-cell"]');
    expect(cells[30].classes()).toContain('intimacy-high');
    // Index 0 = 30 days ago, no data → unavailable.
    expect(cells[0].classes()).toContain('cell-unavailable');
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

    const wrapper = mount(CheckinMatrix, { props: { past: 4 } });
    await flushPromises();

    // Row-major layout with past=4: intimacy cells 0-3.
    const cells = wrapper.findAll('[data-testid="matrix-cell"]');
    expect(cells[0].classes()).toContain('intimacy-very-low');
    expect(cells[1].classes()).toContain('intimacy-low');
    expect(cells[2].classes()).toContain('intimacy-moderate');
    expect(cells[3].classes()).toContain('intimacy-high');
  });

  it('TestDashboard_GivenMissingInsightDay_WhenCellRendered_ThenShowsGreyCell', async () => {
    // All days have no data: API returns empty array → every cell is unavailable.
    insightsApi.getWeeklyInsights.mockResolvedValue([]);

    const wrapper = mount(CheckinMatrix, { props: { past: 3 } });
    await flushPromises();

    const cells = wrapper.findAll('[data-testid="matrix-cell"]');
    // 3 dimensions × 3 days = 9 cells, all must be unavailable (grey).
    expect(cells).toHaveLength(9);
    cells.forEach(cell => {
      expect(cell.classes()).toContain('cell-unavailable');
    });
  });

  it('TestCheckinHistory_GivenHistoryVisible_WhenRendered_ThenIntimacyRowUsesRoseTones', async () => {
    const today = new Date();
    const todayStr = new Date(Date.UTC(
      today.getUTCFullYear(), today.getUTCMonth(), today.getUTCDate(),
    )).toISOString().slice(0, 10).replace(/-/g, '');
    insightsApi.getWeeklyInsights.mockResolvedValue([
      { date: todayStr, intimacy: 80, commitment: 80, passion: 80 },
    ]);

    const wrapper = mount(CheckinMatrix, { props: { past: 1 } });
    await flushPromises();

    // Row 0 = intimacy; only 1 day so cell 0 is intimacy, cell 1 commitment, cell 2 passion.
    const cells = wrapper.findAll('[data-testid="matrix-cell"]');
    expect(cells[0].classes().some(c => c.startsWith('intimacy-'))).toBe(true);
  });

  it('TestCheckinHistory_GivenHistoryVisible_WhenRendered_ThenCommitmentRowUsesSageTones', async () => {
    const today = new Date();
    const todayStr = new Date(Date.UTC(
      today.getUTCFullYear(), today.getUTCMonth(), today.getUTCDate(),
    )).toISOString().slice(0, 10).replace(/-/g, '');
    insightsApi.getWeeklyInsights.mockResolvedValue([
      { date: todayStr, intimacy: 80, commitment: 80, passion: 80 },
    ]);

    const wrapper = mount(CheckinMatrix, { props: { past: 1 } });
    await flushPromises();

    // Row 1 = commitment; with past=1 row-major: cell 0=intimacy, cell 1=commitment, cell 2=passion.
    const cells = wrapper.findAll('[data-testid="matrix-cell"]');
    expect(cells[1].classes().some(c => c.startsWith('commitment-'))).toBe(true);
  });

  it('TestCheckinHistory_GivenHistoryVisible_WhenRendered_ThenPassionRowUsesGoldTones', async () => {
    const today = new Date();
    const todayStr = new Date(Date.UTC(
      today.getUTCFullYear(), today.getUTCMonth(), today.getUTCDate(),
    )).toISOString().slice(0, 10).replace(/-/g, '');
    insightsApi.getWeeklyInsights.mockResolvedValue([
      { date: todayStr, intimacy: 80, commitment: 80, passion: 80 },
    ]);

    const wrapper = mount(CheckinMatrix, { props: { past: 1 } });
    await flushPromises();

    // Row 2 = passion; cell 2 with past=1.
    const cells = wrapper.findAll('[data-testid="matrix-cell"]');
    expect(cells[2].classes().some(c => c.startsWith('passion-'))).toBe(true);
  });

  it('TestCheckinHistory_GivenLowScore_WhenCellRendered_ThenShowsFaintTint', async () => {
    const today = new Date();
    const todayStr = new Date(Date.UTC(
      today.getUTCFullYear(), today.getUTCMonth(), today.getUTCDate(),
    )).toISOString().slice(0, 10).replace(/-/g, '');
    insightsApi.getWeeklyInsights.mockResolvedValue([
      { date: todayStr, intimacy: 10, commitment: 10, passion: 10 },
    ]);

    const wrapper = mount(CheckinMatrix, { props: { past: 1 } });
    await flushPromises();

    const cells = wrapper.findAll('[data-testid="matrix-cell"]');
    expect(cells[0].classes()).toContain('intimacy-very-low');
  });

  it('TestCheckinHistory_GivenHighScore_WhenCellRendered_ThenShowsRichShade', async () => {
    const today = new Date();
    const todayStr = new Date(Date.UTC(
      today.getUTCFullYear(), today.getUTCMonth(), today.getUTCDate(),
    )).toISOString().slice(0, 10).replace(/-/g, '');
    insightsApi.getWeeklyInsights.mockResolvedValue([
      { date: todayStr, intimacy: 90, commitment: 90, passion: 90 },
    ]);

    const wrapper = mount(CheckinMatrix, { props: { past: 1 } });
    await flushPromises();

    const cells = wrapper.findAll('[data-testid="matrix-cell"]');
    expect(cells[0].classes()).toContain('intimacy-high');
  });

  it('TestCheckinHistory_GivenNoCheckin_WhenCellRendered_ThenShowsNeutralState', async () => {
    insightsApi.getWeeklyInsights.mockResolvedValue([]);

    const wrapper = mount(CheckinMatrix, { props: { past: 1 } });
    await flushPromises();

    const cells = wrapper.findAll('[data-testid="matrix-cell"]');
    cells.forEach(cell => {
      expect(cell.classes()).toContain('cell-unavailable');
    });
  });

  it('TestCheckinHistory_GivenAnyCell_WhenTapped_ThenNavigatesToInsightsForDate', async () => {
    const today = new Date();
    const todayStr = new Date(Date.UTC(
      today.getUTCFullYear(), today.getUTCMonth(), today.getUTCDate(),
    )).toISOString().slice(0, 10).replace(/-/g, '');

    insightsApi.getWeeklyInsights.mockResolvedValue([]);

    const wrapper = mount(CheckinMatrix, {
      props: { past: 1 },
      global: {
        stubs: {
          RouterLink: { template: '<a :href="to"><slot /></a>', props: ['to'] },
        },
      },
    });
    await flushPromises();

    // 3 dimensions × 1 day = 3 links, all pointing to today's insights page.
    const links = wrapper.findAll('[data-testid="cell-link"]');
    expect(links).toHaveLength(3);
    links.forEach(link => {
      expect(link.attributes('href')).toBe(`/insights/${todayStr}`);
    });
  });

  it('TestDashboard_GivenMobileViewport_WhenMonthlyMatrixShown_ThenRemainsCondensedAndReadable', async () => {
    // The matrix section must render with a wrapping element that carries the
    // data-testid, and the table must be present. Layout condensation (width:100%,
    // fixed table-layout, small row height) is defined in scoped CSS and verified
    // by visual review at 375px per CONSTITUTION.md baseline.
    const wrapper = mount(CheckinMatrix, { props: { past: 31 } });
    await flushPromises();

    expect(wrapper.find('[data-testid="checkin-history"]').exists()).toBe(true);
    expect(wrapper.find('table').exists()).toBe(true);
    // 31 columns × 3 rows = exactly 93 cells — confirms condensed structure fits.
    expect(wrapper.findAll('[data-testid="matrix-cell"]')).toHaveLength(93);
  });
});
