import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount, flushPromises } from '@vue/test-utils';
import CheckinView from './CheckinView.vue';
import * as checkinApi from '../api/checkin.js';

vi.mock('../api/checkin.js', () => ({
  getTodayCheckin: vi.fn().mockResolvedValue(null),
  saveTodayCheckin: vi.fn().mockResolvedValue({
    felt_close: 7, positive_energy: 8, supported: 6,
    communication_healthy: 9, stress_level: 4, note: ''
  })
}));

const EXISTING_ENTRY = {
  felt_close: 7, positive_energy: 8, supported: 6,
  communication_healthy: 9, stress_level: 4, note: 'good day'
};

describe('CheckinView', () => {
  beforeEach(() => { vi.clearAllMocks(); });

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
