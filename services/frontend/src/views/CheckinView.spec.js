import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount, flushPromises } from '@vue/test-utils';
import CheckinView from './CheckinView.vue';
import * as checkinApi from '../api/checkin.js';

vi.mock('../api/checkin.js', () => ({
  getTodayCheckin: vi.fn().mockResolvedValue(null),
  saveTodayCheckin: vi.fn().mockResolvedValue({
    felt_close: 2, positive_energy: 3, supported: 1,
    communication_healthy: 4, stress_level: -1, note: ''
  })
}));

const EXISTING_ENTRY = {
  felt_close: 2, positive_energy: 3, supported: 1,
  communication_healthy: 4, stress_level: -1, note: 'good day'
};

describe('CheckinView', () => {
  beforeEach(() => { vi.clearAllMocks(); });

  it('TestCheckin_GivenFreshForm_WhenDisplayed_ThenAllSlidersUnsetAtZero', async () => {
    const wrapper = mount(CheckinView);
    await flushPromises();

    const sliders = wrapper.findAll('input[type="range"]');
    expect(sliders).toHaveLength(5);
    for (const slider of sliders) {
      expect(slider.attributes('data-unset')).toBe('true');
      expect(slider.element.value).toBe('0');
    }
  });


  it('TestCheckin_GivenUnsetSlider_WhenUserInteracts_ThenUnsetStyleClears', async () => {
    const wrapper = mount(CheckinView);
    await flushPromises();

    const slider = wrapper.find('input[type="range"]');
    expect(slider.attributes('data-unset')).toBe('true');

    await slider.setValue(1);

    expect(slider.attributes('data-unset')).toBeUndefined();
    expect(slider.classes()).not.toContain('checkin-slider--unset');
  });

  it('TestCheckin_GivenFormRendered_WhenDisplayed_ThenEachSliderHasDescriptionText', () => {
    const wrapper = mount(CheckinView);

    const descriptions = [
      'Did you feel emotionally connected?',
      'Was there lightness or joy between you?',
      'Did you feel like you had each other\'s backs?',
      'Were you able to talk openly and be heard?',
      'How much did personal stress weigh on you?',
    ];
    for (const text of descriptions) {
      expect(wrapper.text()).toContain(text);
    }
  });

  it('TestCheckin_GivenFormRendered_WhenDisplayed_ThenSliderRangeIsNegFiveToPosFive', () => {
    const wrapper = mount(CheckinView);

    const sliders = wrapper.findAll('input[type="range"]');
    expect(sliders).toHaveLength(5);
    for (const slider of sliders) {
      expect(slider.attributes('min')).toBe('-5');
      expect(slider.attributes('max')).toBe('5');
    }
  });

  it('TestCheckin_GivenCheckinPage_WhenLoaded_ThenFiveSlidersVisible', () => {
    const wrapper = mount(CheckinView);

    const labels = [
      'Felt close today',
      'Positive energy / fun',
      'Supported / team',
      'Communication healthy',
      'My stress level',
    ];
    for (const label of labels) {
      expect(wrapper.text()).toContain(label);
    }

    const sliders = wrapper.findAll('input[type="range"]');
    expect(sliders).toHaveLength(5);
  });

  it('TestCheckin_GivenNoExistingEntry_WhenPageLoads_ThenSlidersShowUnsetState', async () => {
    const wrapper = mount(CheckinView);
    await flushPromises();

    const sliders = wrapper.findAll('input[type="range"]');
    for (const slider of sliders) {
      expect(slider.attributes('data-unset')).toBe('true');
    }
  });

  it('TestCheckin_GivenValuesSet_WhenSaved_ThenConfirmationShown', async () => {
    const wrapper = mount(CheckinView);
    await flushPromises();

    await wrapper.find('button[data-testid="save-checkin"]').trigger('click');
    await flushPromises();

    expect(wrapper.find('[data-testid="checkin-confirmation"]').exists()).toBe(true);
  });

  it('TestCheckin_GivenExistingEntry_WhenPageLoads_ThenFormPrefilled', async () => {
    checkinApi.getTodayCheckin.mockResolvedValue(EXISTING_ENTRY);
    const wrapper = mount(CheckinView);
    await flushPromises();

    const sliders = wrapper.findAll('input[type="range"]');
    for (const slider of sliders) {
      expect(slider.attributes('data-unset')).toBeUndefined();
    }
    expect(wrapper.find('textarea').element.value).toBe('good day');
  });

  it('TestCheckin_GivenExistingEntry_WhenResubmitted_ThenEntryUpdated', async () => {
    checkinApi.getTodayCheckin.mockResolvedValue(EXISTING_ENTRY);
    checkinApi.saveTodayCheckin.mockResolvedValue({ ...EXISTING_ENTRY, felt_close: 9 });

    const wrapper = mount(CheckinView);
    await flushPromises();

    await wrapper.find('button[data-testid="save-checkin"]').trigger('click');
    await flushPromises();

    expect(checkinApi.saveTodayCheckin).toHaveBeenCalledOnce();
    expect(wrapper.find('[data-testid="checkin-confirmation"]').exists()).toBe(true);
  });
});
