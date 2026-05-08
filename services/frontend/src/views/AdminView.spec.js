import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount, flushPromises } from '@vue/test-utils';
import AdminView from './AdminView.vue';
import * as adminApi from '../api/admin.js';

vi.mock('vue-router', () => ({
  useRouter: vi.fn().mockReturnValue({ push: vi.fn() }),
}));

vi.mock('../api/admin.js', () => ({
  getAdminUsers: vi.fn(),
  activateUser: vi.fn(),
  deactivateUser: vi.fn(),
}));

const stubs = { NavBar: true };

describe('AdminView', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    adminApi.getAdminUsers.mockResolvedValue([]);
    adminApi.activateUser.mockResolvedValue(undefined);
    adminApi.deactivateUser.mockResolvedValue(undefined);
  });

  it('TestAdminView_GivenUsers_WhenRendered_ThenAllUsersListed', async () => {
    adminApi.getAdminUsers.mockResolvedValue([
      { id: '1', email: 'a@test.com', firstName: 'Alice', role: 'user', isActive: true, createdAt: '2026-01-01T00:00:00Z' },
    ]);
    const wrapper = mount(AdminView, { global: { stubs } });
    await flushPromises();
    expect(wrapper.text()).toContain('a@test.com');
  });

  it('TestAdminView_GivenActiveUser_WhenDeactivateClicked_ThenDeactivateApiCalled', async () => {
    adminApi.getAdminUsers.mockResolvedValue([
      { id: '1', email: 'a@test.com', firstName: 'Alice', role: 'user', isActive: true, createdAt: '2026-01-01T00:00:00Z' },
    ]);
    const wrapper = mount(AdminView, { global: { stubs } });
    await flushPromises();
    await wrapper.find('[data-testid="deactivate-1"]').trigger('click');
    expect(adminApi.deactivateUser).toHaveBeenCalledWith('1');
  });

  it('TestAdminView_GivenInactiveUser_WhenActivateClicked_ThenActivateApiCalled', async () => {
    adminApi.getAdminUsers.mockResolvedValue([
      { id: '1', email: 'a@test.com', firstName: 'Alice', role: 'user', isActive: false, createdAt: '2026-01-01T00:00:00Z' },
    ]);
    const wrapper = mount(AdminView, { global: { stubs } });
    await flushPromises();
    await wrapper.find('[data-testid="activate-1"]').trigger('click');
    expect(adminApi.activateUser).toHaveBeenCalledWith('1');
  });

  it('TestAdminView_GivenLoadingState_WhenFetching_ThenSpinnerVisible', async () => {
    let resolve;
    adminApi.getAdminUsers.mockReturnValue(new Promise(r => { resolve = r; }));
    const wrapper = mount(AdminView, { global: { stubs } });
    expect(wrapper.find('[data-testid="loading"]').exists()).toBe(true);
    resolve([]);
    await flushPromises();
    expect(wrapper.find('[data-testid="loading"]').exists()).toBe(false);
  });
});
