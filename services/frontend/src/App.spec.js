import { describe, it, expect, vi } from 'vitest';
import { mount } from '@vue/test-utils';
import App from './App.vue';

vi.mock('vue-router', () => ({
  useRoute: vi.fn(() => ({ path: '/login' })),
  RouterView: { template: '<div />' },
}));

describe('App', () => {
  it('TestApp_GivenLoginRoute_WhenRendered_ThenNavBarAbsent', () => {
    const wrapper = mount(App, {
      global: {
        stubs: {
          NavBar: { template: '<div data-testid="navbar" />' },
          RouterView: { template: '<div />' },
        },
      },
    });

    expect(wrapper.find('[data-testid="navbar"]').exists()).toBe(false);
  });
});
