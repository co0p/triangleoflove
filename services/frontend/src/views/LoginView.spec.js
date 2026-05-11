import { describe, it, expect, vi, beforeEach } from 'vitest';
import { RouterLinkStub, mount } from '@vue/test-utils';
import { useRoute, useRouter } from 'vue-router';
import LoginView from './LoginView.vue';

vi.mock('vue-router', () => ({
  useRoute: vi.fn().mockReturnValue({ query: {} }),
  useRouter: vi.fn().mockReturnValue({ push: vi.fn() }),
}));

describe('LoginView', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    useRoute.mockReturnValue({ query: {} });
    useRouter.mockReturnValue({ push: vi.fn() });
  });

  it('TestRegister_GivenLoginPage_WhenSignupChosen_ThenRegistrationOpens', () => {
    const wrapper = mount(LoginView, {
      global: {
        stubs: {
          RouterLink: RouterLinkStub,
        },
      },
    });

    const signupLink = wrapper.findComponent(RouterLinkStub);

    expect(signupLink.exists()).toBe(true);
    expect(signupLink.props('to')).toBe('/register');
    expect(wrapper.text()).toContain('Sign up');
  });

  it('TestRegister_GivenValidRegistration_WhenSucceeded_ThenLoginWithSuccessShown', () => {
    useRoute.mockReturnValue({ query: { registered: '1' } });
    const wrapper = mount(LoginView, {
      global: {
        stubs: {
          RouterLink: RouterLinkStub,
        },
      },
    });

    expect(wrapper.find('[role="status"]').text()).toContain('Account created. You can now sign in.');
  });
});