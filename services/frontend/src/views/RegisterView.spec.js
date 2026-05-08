import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount, flushPromises } from '@vue/test-utils';
import { useRouter } from 'vue-router';
import RegisterView from './RegisterView.vue';
import * as authApi from '../api/auth.js';

vi.mock('vue-router', () => ({
  useRouter: vi.fn().mockReturnValue({ push: vi.fn() }),
}));

vi.mock('../api/auth.js', () => ({
  register: vi.fn(),
}));

const stubs = {
  'router-link': true,
};

describe('RegisterView', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    useRouter.mockReturnValue({ push: vi.fn() });
  });

  it('TestRegistration_GivenValidInput_WhenSubmitted_ThenRedirectsToLogin', async () => {
    authApi.register.mockResolvedValue(undefined);
    const wrapper = mount(RegisterView, { global: { stubs } });

    await wrapper.find('input[type="email"]').setValue('new@test.com');
    await wrapper.find('input[type="password"]').setValue('securepass');
    await wrapper.find('input[name="firstName"]').setValue('Alice');
    await wrapper.find('form').trigger('submit');
    await flushPromises();

    expect(authApi.register).toHaveBeenCalledWith('new@test.com', 'securepass', 'Alice');
    expect(useRouter().push).toHaveBeenCalledWith('/login');
  });

  it('TestRegistration_GivenDuplicateEmail_WhenSubmitted_ThenErrorShown', async () => {
    authApi.register.mockRejectedValue(new Error('duplicate email'));
    const wrapper = mount(RegisterView, { global: { stubs } });

    await wrapper.find('input[type="email"]').setValue('existing@test.com');
    await wrapper.find('input[type="password"]').setValue('securepass');
    await wrapper.find('input[name="firstName"]').setValue('Alice');
    await wrapper.find('form').trigger('submit');
    await flushPromises();

    expect(wrapper.find('[role="alert"]').exists()).toBe(true);
  });

  it('TestRegistration_GivenShortPassword_WhenSubmitted_ThenErrorShown', async () => {
    const wrapper = mount(RegisterView, { global: { stubs } });

    await wrapper.find('input[type="email"]').setValue('new@test.com');
    await wrapper.find('input[type="password"]').setValue('short');
    await wrapper.find('input[name="firstName"]').setValue('Alice');
    await wrapper.find('form').trigger('submit');
    await flushPromises();

    expect(authApi.register).not.toHaveBeenCalled();
    expect(wrapper.find('[role="alert"]').exists()).toBe(true);
  });
});
