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

  it('TestRegister_GivenEmailInput_WhenFormatValid_ThenEmailCheckShown', async () => {
    const wrapper = mount(RegisterView, { global: { stubs } });

    await wrapper.find('input[type="email"]').setValue('new@test.com');

    const emailIndicator = wrapper.get('[data-testid="email-format-indicator"]');

    expect(emailIndicator.classes()).toContain('validation-rule--met');
    expect(emailIndicator.text()).toContain('Ready');
  });

  it('TestRegister_GivenPassword_WhenMinLengthMet_ThenLengthCheckShown', async () => {
    const wrapper = mount(RegisterView, { global: { stubs } });

    await wrapper.find('input[name="password"]').setValue('securepass');

    const lengthIndicator = wrapper.get('[data-testid="password-length-indicator"]');

    expect(lengthIndicator.classes()).toContain('validation-rule--met');
    expect(lengthIndicator.text()).toContain('Ready');
  });

  it('TestRegister_GivenPassword_WhenSpecialCharPresent_ThenCharCheckShown', async () => {
    const wrapper = mount(RegisterView, { global: { stubs } });

    await wrapper.find('input[name="password"]').setValue('securepass!');

    const specialIndicator = wrapper.get('[data-testid="password-special-indicator"]');

    expect(specialIndicator.classes()).toContain('validation-rule--met');
    expect(specialIndicator.text()).toContain('Ready');
  });

  it('TestRegister_GivenTwoPasswords_WhenTheyMatch_ThenMatchCheckShown', async () => {
    const wrapper = mount(RegisterView, { global: { stubs } });

    await wrapper.find('input[name="password"]').setValue('securepass!');
    await wrapper.find('input[name="passwordConfirmation"]').setValue('securepass!');

    const matchIndicator = wrapper.get('[data-testid="password-match-indicator"]');

    expect(matchIndicator.classes()).toContain('validation-rule--met');
    expect(matchIndicator.text()).toContain('Ready');
  });

  it('TestRegister_GivenRegistrationForm_WhenCompleted_ThenTwoPasswordsRequired', async () => {
    const wrapper = mount(RegisterView, { global: { stubs } });

    await wrapper.find('input[type="email"]').setValue('new@test.com');
    await wrapper.find('input[name="password"]').setValue('securepass');
    await wrapper.find('input[name="firstName"]').setValue('Alice');
    await wrapper.find('form').trigger('submit');
    await flushPromises();

    expect(authApi.register).not.toHaveBeenCalled();
    expect(wrapper.find('[role="alert"]').text()).toContain('Please confirm your password.');
  });

  it('TestRegister_GivenValidRegistration_WhenSucceeded_ThenLoginWithSuccessShown', async () => {
    authApi.register.mockResolvedValue(undefined);
    const wrapper = mount(RegisterView, { global: { stubs } });

    await wrapper.find('input[type="email"]').setValue('new@test.com');
    await wrapper.find('input[name="password"]').setValue('securepass!');
    await wrapper.find('input[name="passwordConfirmation"]').setValue('securepass!');
    await wrapper.find('input[name="firstName"]').setValue('Alice');
    await wrapper.find('form').trigger('submit');
    await flushPromises();

    expect(authApi.register).toHaveBeenCalledWith('new@test.com', 'securepass!', 'Alice');
    expect(useRouter().push).toHaveBeenCalledWith({ path: '/login', query: { registered: '1' } });
  });

  it('TestRegister_GivenBadEmail_WhenValidated_ThenEmailErrorShown', async () => {
    const wrapper = mount(RegisterView, { global: { stubs } });

    await wrapper.find('input[type="email"]').setValue('not-an-email');
    await wrapper.find('input[name="password"]').setValue('securepass!');
    await wrapper.find('input[name="passwordConfirmation"]').setValue('securepass!');
    await wrapper.find('input[name="firstName"]').setValue('Alice');
    await wrapper.find('form').trigger('submit');
    await flushPromises();

    expect(authApi.register).not.toHaveBeenCalled();
    expect(wrapper.find('[role="alert"]').text()).toContain('Enter a valid email address.');
  });

  it('TestRegister_GivenPasswordMismatch_WhenValidated_ThenMismatchErrorShown', async () => {
    const wrapper = mount(RegisterView, { global: { stubs } });

    await wrapper.find('input[type="email"]').setValue('new@test.com');
    await wrapper.find('input[name="password"]').setValue('securepass!');
    await wrapper.find('input[name="passwordConfirmation"]').setValue('otherpass!');
    await wrapper.find('input[name="firstName"]').setValue('Alice');
    await wrapper.find('form').trigger('submit');
    await flushPromises();

    expect(authApi.register).not.toHaveBeenCalled();
    expect(wrapper.find('[role="alert"]').text()).toContain('Passwords must match.');
  });

  it('TestRegister_GivenShortPassword_WhenValidated_ThenLengthErrorShown', async () => {
    const wrapper = mount(RegisterView, { global: { stubs } });

    await wrapper.find('input[type="email"]').setValue('new@test.com');
    await wrapper.find('input[name="password"]').setValue('short!');
    await wrapper.find('input[name="passwordConfirmation"]').setValue('short!');
    await wrapper.find('input[name="firstName"]').setValue('Alice');
    await wrapper.find('form').trigger('submit');
    await flushPromises();

    expect(authApi.register).not.toHaveBeenCalled();
    expect(wrapper.find('[role="alert"]').text()).toContain('Password must be at least 8 characters.');
  });

  it('TestRegister_GivenNoSpecialChar_WhenValidated_ThenCharErrorShown', async () => {
    const wrapper = mount(RegisterView, { global: { stubs } });

    await wrapper.find('input[type="email"]').setValue('new@test.com');
    await wrapper.find('input[name="password"]').setValue('securepass');
    await wrapper.find('input[name="passwordConfirmation"]').setValue('securepass');
    await wrapper.find('input[name="firstName"]').setValue('Alice');
    await wrapper.find('form').trigger('submit');
    await flushPromises();

    expect(authApi.register).not.toHaveBeenCalled();
    expect(wrapper.find('[role="alert"]').text()).toContain('Password must include a special character.');
  });

  it('TestRegister_GivenInvalidRules_WhenSubmitted_ThenAccountCreationRejected', async () => {
    authApi.register.mockRejectedValue(new Error('Password must include a special character.'));
    const wrapper = mount(RegisterView, { global: { stubs } });

    await wrapper.find('input[type="email"]').setValue('new@test.com');
    await wrapper.find('input[name="password"]').setValue('securepass!');
    await wrapper.find('input[name="passwordConfirmation"]').setValue('securepass!');
    await wrapper.find('input[name="firstName"]').setValue('Alice');
    await wrapper.find('form').trigger('submit');
    await flushPromises();

    expect(wrapper.find('[role="alert"]').text()).toContain('Password must include a special character.');
    expect(useRouter().push).not.toHaveBeenCalled();
  });

  it('TestRegistration_GivenDuplicateEmail_WhenSubmitted_ThenErrorShown', async () => {
    authApi.register.mockRejectedValue(new Error('duplicate email'));
    const wrapper = mount(RegisterView, { global: { stubs } });

    await wrapper.find('input[type="email"]').setValue('existing@test.com');
    await wrapper.find('input[name="password"]').setValue('securepass!');
    await wrapper.find('input[name="passwordConfirmation"]').setValue('securepass!');
    await wrapper.find('input[name="firstName"]').setValue('Alice');
    await wrapper.find('form').trigger('submit');
    await flushPromises();

    expect(wrapper.find('[role="alert"]').exists()).toBe(true);
  });
});
