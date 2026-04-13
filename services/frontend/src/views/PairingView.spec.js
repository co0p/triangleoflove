import { describe, it, expect, vi, beforeEach } from 'vitest';
import { mount, flushPromises } from '@vue/test-utils';
import PairingView from './PairingView.vue';
import * as pairingApi from '../api/pairing.js';

vi.mock('vue-router', () => ({
  useRouter: vi.fn().mockReturnValue({ push: vi.fn() }),
}));

vi.mock('../api/pairing.js', () => ({
  getCoupleStatus: vi.fn(),
  getPairing: vi.fn(),
  regeneratePairing: vi.fn(),
  connectPairing: vi.fn(),
  unpairCouple: vi.fn(),
}));

const stubs = { NavBar: true };

const PAIRED_STATUS = {
  paired: true,
  partner_first_name: 'Jordan',
  paired_since: '2026-04-01T00:00:00Z',
};

const UNPAIRED_STATUS = { paired: false };

describe('PairingView', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    pairingApi.getPairing.mockResolvedValue({ invite_code: 'ABC123' });
  });

  it('TestPairingView_GivenPaired_WhenRendered_ThenUnpairButtonVisibleAndCodeHidden', async () => {
    pairingApi.getCoupleStatus.mockResolvedValue(PAIRED_STATUS);

    const wrapper = mount(PairingView, { global: { stubs } });
    await flushPromises();

    expect(wrapper.find('[data-testid="unpair-button"]').exists()).toBe(true);
    expect(wrapper.find('[data-testid="invite-code"]').exists()).toBe(false);
    expect(wrapper.find('[data-testid="partner-code-input"]').exists()).toBe(false);
  });

  it('TestPairingView_GivenPaired_WhenUnpairClicked_ThenConfirmationModalAppears', async () => {
    pairingApi.getCoupleStatus.mockResolvedValue(PAIRED_STATUS);

    const wrapper = mount(PairingView, { global: { stubs } });
    await flushPromises();

    await wrapper.find('[data-testid="unpair-button"]').trigger('click');

    expect(wrapper.find('[data-testid="unpair-modal"]').exists()).toBe(true);
    expect(wrapper.find('[data-testid="unpair-modal-confirm"]').exists()).toBe(true);
    expect(wrapper.find('[data-testid="unpair-modal-cancel"]').exists()).toBe(true);
  });

  it('TestPairingView_GivenModalOpen_WhenCancelSelected_ThenRemainsPaired', async () => {
    pairingApi.getCoupleStatus.mockResolvedValue(PAIRED_STATUS);

    const wrapper = mount(PairingView, { global: { stubs } });
    await flushPromises();

    await wrapper.find('[data-testid="unpair-button"]').trigger('click');
    await wrapper.find('[data-testid="unpair-modal-cancel"]').trigger('click');

    expect(wrapper.find('[data-testid="unpair-modal"]').exists()).toBe(false);
    expect(wrapper.find('[data-testid="partner-name"]').exists()).toBe(true);
    expect(wrapper.find('[data-testid="unpair-button"]').exists()).toBe(true);
    expect(pairingApi.unpairCouple).not.toHaveBeenCalled();
  });

  it('TestPairingView_GivenModalOpen_WhenYesSelected_ThenBothUsersUnpaired', async () => {
    pairingApi.getCoupleStatus.mockResolvedValue(PAIRED_STATUS);
    pairingApi.unpairCouple.mockResolvedValue({ status: 'unpaired' });

    const wrapper = mount(PairingView, { global: { stubs } });
    await flushPromises();

    await wrapper.find('[data-testid="unpair-button"]').trigger('click');
    await wrapper.find('[data-testid="unpair-modal-confirm"]').trigger('click');
    await flushPromises();

    expect(pairingApi.unpairCouple).toHaveBeenCalledOnce();
    expect(wrapper.find('[data-testid="unpair-modal"]').exists()).toBe(false);
    expect(wrapper.find('[data-testid="partner-name"]').exists()).toBe(false);
  });

  it('TestPairingView_GivenJustUnpaired_WhenViewUpdates_ThenInviteCodeVisible', async () => {
    pairingApi.getCoupleStatus.mockResolvedValue(PAIRED_STATUS);
    pairingApi.unpairCouple.mockResolvedValue({ status: 'unpaired' });

    const wrapper = mount(PairingView, { global: { stubs } });
    await flushPromises();

    await wrapper.find('[data-testid="unpair-button"]').trigger('click');
    await wrapper.find('[data-testid="unpair-modal-confirm"]').trigger('click');
    await flushPromises();

    expect(wrapper.find('[data-testid="invite-code"]').exists()).toBe(true);
    expect(wrapper.find('[data-testid="invite-code"]').text()).toBe('ABC123');
    expect(wrapper.find('[data-testid="unpair-button"]').exists()).toBe(false);
  });
});
