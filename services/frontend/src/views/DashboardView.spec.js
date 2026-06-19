import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount, flushPromises } from '@vue/test-utils';
import DashboardView from './DashboardView.vue';
import * as usersApi from '../api/users.js';
import * as pairingApi from '../api/pairing.js';
import * as sessionApi from '../api/checkin.js';
import * as insightsApi from '../api/insights.js';
import { useCurrentUser } from '../composables/useCurrentUser.js';

const { push } = vi.hoisted(() => ({
  push: vi.fn(),
}));

vi.mock('vue-router', () => ({
  useRouter: vi.fn().mockReturnValue({ push }),
}));

vi.mock('../api/users.js', () => ({
  getMe: vi.fn(),
}));

vi.mock('../api/pairing.js', () => ({
  getCoupleStatus: vi.fn(),
}));

vi.mock('../api/checkin.js', () => ({
  getTodaySession: vi.fn(),
}));

vi.mock('../api/insights.js', () => ({
  getWeeklyInsights: vi.fn(),
}));

const stubs = {
  NavBar: true,
  RouterLink: { template: '<a v-bind="$attrs"><slot /></a>' },
  CheckinMatrix: { template: '<div data-testid="checkin-history" />' },
};

describe('DashboardView', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.setItem('token', 'test-token');
    useCurrentUser().reset();
    usersApi.getMe.mockResolvedValue({ firstName: 'Alice' });
    pairingApi.getCoupleStatus.mockResolvedValue({ paired: false, partner_first_name: '' });
    sessionApi.getTodaySession.mockResolvedValue(null);
    insightsApi.getWeeklyInsights.mockResolvedValue([{ date: '20260619', intimacy: 60, commitment: 60, passion: 60 }]);
  });

  it('TestDashboard_GivenPartnerName_WhenRendered_ThenPairingStatusVisible', async () => {
    pairingApi.getCoupleStatus.mockResolvedValue({ paired: true, partner_first_name: 'Bob' });

    const wrapper = mount(DashboardView, { global: { stubs } });
    await flushPromises();

    expect(wrapper.find('[data-testid="pairing-status"]').exists()).toBe(true);
    expect(wrapper.find('[data-testid="pairing-status"]').text()).toContain('Bob');
  });

  it('TestDashboard_GivenNoPartner_WhenRendered_ThenNotConnectedVisible', async () => {
    const wrapper = mount(DashboardView, { global: { stubs } });
    await flushPromises();

    expect(wrapper.text()).toContain('Not connected yet');
  });

  it('TestDashboard_GivenSessionExists_WhenRendered_ThenSessionCompletedStatusVisible', async () => {
    sessionApi.getTodaySession.mockResolvedValue({ felt_close: 3 });

    const wrapper = mount(DashboardView, { global: { stubs } });
    await flushPromises();

    expect(wrapper.text()).toContain('Weekly check-in logged today');
  });

  it('TestDashboard_GivenRecentCheckin_WhenRendered_ThenBacklogShowsOneDot', async () => {
    insightsApi.getWeeklyInsights.mockResolvedValue([{ date: '20260619', intimacy: 60, commitment: 60, passion: 60 }]);

    const wrapper = mount(DashboardView, { global: { stubs } });
    await flushPromises();

    expect(wrapper.findAll('.weekly-backlog-dot--filled')).toHaveLength(1);
  });

  it('TestDashboard_GivenThreeDaysBehind_WhenRendered_ThenBacklogShowsTwoDots', async () => {
    insightsApi.getWeeklyInsights.mockResolvedValue([{ date: '20260616', intimacy: 60, commitment: 60, passion: 60 }]);

    const wrapper = mount(DashboardView, { global: { stubs } });
    await flushPromises();

    expect(wrapper.findAll('.weekly-backlog-dot--filled')).toHaveLength(2);
  });

  it('TestDashboard_GivenOldCheckin_WhenRendered_ThenBacklogShowsThreeDots', async () => {
    insightsApi.getWeeklyInsights.mockResolvedValue([{ date: '20260610', intimacy: 60, commitment: 60, passion: 60 }]);

    const wrapper = mount(DashboardView, { global: { stubs } });
    await flushPromises();

    expect(wrapper.findAll('.weekly-backlog-dot--filled')).toHaveLength(3);
  });

  it('TestDashboard_GivenPrimaryNavigation_WhenInsightsLinkClicked_ThenOpensWeeklyInsights', async () => {
    const wrapper = mount(DashboardView, { global: { stubs } });
    await flushPromises();

    const insightsLink = wrapper.find('[data-testid="insights-link"]');
    expect(insightsLink.exists()).toBe(true);
    expect(insightsLink.attributes('to')).toBe('/insights');
  });

  it('TestDashboard_GivenUnauthorizedLoad_WhenRendered_ThenClearsTokenAndRedirectsToLogin', async () => {
    const removeItemSpy = vi.spyOn(Storage.prototype, 'removeItem');
    pairingApi.getCoupleStatus.mockRejectedValue(new Error('unauthorized'));

    mount(DashboardView, { global: { stubs } });
    await flushPromises();

    expect(removeItemSpy).toHaveBeenCalledWith('token');
    expect(push).toHaveBeenCalledWith('/login');

    removeItemSpy.mockRestore();
  });

  it('TestDashboard_GivenNonAuthLoadError_WhenRendered_ThenSessionRemainsIntact', async () => {
    const removeItemSpy = vi.spyOn(Storage.prototype, 'removeItem');
    pairingApi.getCoupleStatus.mockRejectedValue(new Error('failed to load couple status'));

    mount(DashboardView, { global: { stubs } });
    await flushPromises();

    expect(removeItemSpy).not.toHaveBeenCalled();
    expect(push).not.toHaveBeenCalled();

    removeItemSpy.mockRestore();
  });

  it('TestDashboard_GivenMonthlyInsights_WhenDashboardLoads_ThenInsightsSectionVisible', async () => {
    const wrapper = mount(DashboardView, { global: { stubs } });
    await flushPromises();

    expect(wrapper.find('[data-testid="checkin-history"]').exists()).toBe(true);
  });

  it('TestDashboard_GivenSharedUserState_WhenRendered_ThenWelcomeHeadingMatchesNavBar', async () => {
    const wrapper = mount(DashboardView, { global: { stubs } });
    await flushPromises();

    expect(wrapper.text()).toContain('Welcome back, Alice');
  });
});
