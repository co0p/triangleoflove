import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount, flushPromises } from '@vue/test-utils';
import NavBar from './NavBar.vue';
import * as usersApi from '../api/users.js';
import { useCurrentUser } from '../composables/useCurrentUser.js';

vi.mock('../api/users.js', () => ({
  getMe: vi.fn(),
}));

const stubs = {
  RouterLink: { template: '<a v-bind="$attrs"><slot /></a>' },
};

describe('NavBar', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.setItem('token', 'test-token');
    useCurrentUser().reset();
  });

  it('TestNavBar_GivenUserStateLoaded_WhenMounted_ThenGreetingVisible', async () => {
    // Simulate state pre-populated before NavBar mounts (as LoginView does)
    usersApi.getMe.mockResolvedValue({ firstName: 'Alice' });
    await useCurrentUser().load();

    const wrapper = mount(NavBar, { global: { stubs } });

    expect(wrapper.text()).toContain('Hello, Alice');
  });

  it('TestNavBar_GivenUserStateCached_WhenMounted_ThenGetMeCalledOnce', async () => {
    usersApi.getMe.mockResolvedValue({ firstName: 'Alice' });

    // First mount — initial load
    const wrapper1 = mount(NavBar, { global: { stubs } });
    await flushPromises();
    wrapper1.unmount();

    // Second mount — state should be cached, no new getMe call
    const wrapper2 = mount(NavBar, { global: { stubs } });
    await flushPromises();

    expect(usersApi.getMe).toHaveBeenCalledTimes(1);
  });
});
