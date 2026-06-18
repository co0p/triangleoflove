import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount, flushPromises } from '@vue/test-utils';
import CheckinView from './CheckinView.vue';
import * as sessionApi from '../api/checkin.js';

vi.mock('../api/checkin.js', () => ({
  getTodaySession: vi.fn().mockResolvedValue(null),
  saveTodaySession: vi.fn().mockResolvedValue({
    felt_understood: 4, meaningful_sharing: 3, could_count_on_them: 5,
    effort_for_us: 2, desire: 3, spark: 4, mood: 4, note: ''
  })
}));

const EXISTING_ENTRY = {
  felt_understood: 4,
  meaningful_sharing: 3,
  could_count_on_them: 5,
  effort_for_us: 2,
  desire: 3,
  spark: 4,
  mood: 4,
  note: 'good day'
};

describe('CheckinView', () => {
  beforeEach(() => {
    vi.resetAllMocks();
    sessionApi.getTodaySession.mockResolvedValue(null);
    sessionApi.saveTodaySession.mockResolvedValue({
      felt_understood: 4, meaningful_sharing: 3, could_count_on_them: 5,
      effort_for_us: 2, desire: 3, spark: 4, mood: 4, note: ''
    });
  });

  // --- Acceptance tests (D2) ---

  it('TestSessionView_GivenPageLoaded_WhenRendered_ThenShowsNewMetricsOnly', async () => {
    const wrapper = mount(CheckinView);
    await flushPromises();

    expect(wrapper.text()).toContain('Felt understood');
    expect(wrapper.text()).toContain('Meaningful sharing');
    expect(wrapper.text()).toContain('Could count on them');
    expect(wrapper.text()).toContain('Effort for us');
    expect(wrapper.text()).toContain('Desire');
    expect(wrapper.text()).toContain('Spark');
    expect(wrapper.text()).toContain('My mood today');

    expect(wrapper.text()).not.toContain('Felt close today');
    expect(wrapper.text()).not.toContain('Positive energy / fun');
    expect(wrapper.text()).not.toContain('Supported / team');
    expect(wrapper.text()).not.toContain('Communication healthy');
    expect(wrapper.text()).not.toContain('My stress level');

    expect(wrapper.findAll('input[type="range"]')).toHaveLength(7);
  });

  it('TestSessionView_GivenAnySlider_WhenRendered_ThenRangeIsOneToFive', async () => {
    const wrapper = mount(CheckinView);
    await flushPromises();

    const sliders = wrapper.findAll('input[type="range"]');
    expect(sliders).toHaveLength(7);
    for (const slider of sliders) {
      expect(slider.attributes('min')).toBe('1');
      expect(slider.attributes('max')).toBe('5');
    }
  });

  it('TestSessionView_GivenSavedEntry_WhenPageReloaded_ThenValuesPrePopulated', async () => {
    sessionApi.getTodaySession.mockResolvedValue(EXISTING_ENTRY);
    const wrapper = mount(CheckinView);
    await flushPromises();

    const sliders = wrapper.findAll('input[type="range"]');
    expect(sliders).toHaveLength(7);
    for (const slider of sliders) {
      expect(slider.attributes('data-unset')).toBeUndefined();
    }
    expect(wrapper.find('textarea').element.value).toBe('good day');
  });

  // --- Behavioural tests ---

  it('TestSession_GivenFreshForm_WhenDisplayed_ThenAllSlidersUnset', async () => {
    const wrapper = mount(CheckinView);
    await flushPromises();

    const sliders = wrapper.findAll('input[type="range"]');
    expect(sliders).toHaveLength(7);
    for (const slider of sliders) {
      expect(slider.attributes('data-unset')).toBe('true');
    }
  });

  it('TestSession_GivenUnsetSlider_WhenUserInteracts_ThenUnsetStyleClears', async () => {
    const wrapper = mount(CheckinView);
    await flushPromises();

    const slider = wrapper.find('input[type="range"]');
    expect(slider.attributes('data-unset')).toBe('true');

    await slider.setValue(3);

    expect(slider.attributes('data-unset')).toBeUndefined();
    expect(slider.classes()).not.toContain('checkin-slider--unset');
  });

  it('TestSession_GivenValuesSet_WhenSaved_ThenConfirmationShown', async () => {
    const wrapper = mount(CheckinView);
    await flushPromises();

    await wrapper.find('button[data-testid="save-checkin"]').trigger('click');
    await flushPromises();

    expect(wrapper.find('[data-testid="checkin-confirmation"]').exists()).toBe(true);
  });

  it('TestSession_GivenExistingEntry_WhenResubmitted_ThenEntryUpdated', async () => {
    sessionApi.getTodaySession.mockResolvedValue(EXISTING_ENTRY);
    sessionApi.saveTodaySession.mockResolvedValue({ ...EXISTING_ENTRY, felt_understood: 5 });

    const wrapper = mount(CheckinView);
    await flushPromises();

    await wrapper.find('button[data-testid="save-checkin"]').trigger('click');
    await flushPromises();

    expect(sessionApi.saveTodaySession).toHaveBeenCalledOnce();
    expect(wrapper.find('[data-testid="checkin-confirmation"]').exists()).toBe(true);
  });
});
