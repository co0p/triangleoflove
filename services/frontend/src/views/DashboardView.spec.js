import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount, flushPromises } from '@vue/test-utils';
import DashboardView from './DashboardView.vue';
import * as usersApi from '../api/users.js';
import * as pairingApi from '../api/pairing.js';
import * as checkinApi from '../api/checkin.js';

vi.mock('vue-router', () => ({
  useRouter: vi.fn().mockReturnValue({ push: vi.fn() }),
}));

vi.mock('../api/users.js', () => ({
  getMe: vi.fn(),
}));

vi.mock('../api/pairing.js', () => ({
  getCoupleStatus: vi.fn(),
}));

vi.mock('../api/checkin.js', () => ({
  getTodayCheckin: vi.fn(),
}));

const stubs = {
  NavBar: true,
  RouterLink: { template: '<a><slot /></a>' },
};

describe('DashboardView', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    usersApi.getMe.mockResolvedValue({ firstName: 'Alice' });
    pairingApi.getCoupleStatus.mockResolvedValue({ paired: false, partner_first_name: '' });
    checkinApi.getTodayCheckin.mockResolvedValue(null);
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

  it('TestDashboard_GivenCheckinExists_WhenRendered_ThenCheckedInStatusVisible', async () => {
    checkinApi.getTodayCheckin.mockResolvedValue({ felt_close: 3 });

    const wrapper = mount(DashboardView, { global: { stubs } });
    await flushPromises();

    expect(wrapper.text()).toContain('Checked in today');
  });
});
