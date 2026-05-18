import { describe, it, expect, vi, beforeEach } from 'vitest';
import * as usersApi from '../api/users.js';
import { useCurrentUser } from './useCurrentUser.js';

vi.mock('../api/users.js', () => ({
  getMe: vi.fn(),
}));

describe('useCurrentUser', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    localStorage.removeItem('token');
    useCurrentUser().reset();
  });

  it('TestCurrentUser_GivenTokenChanges_WhenLoadCalledAgain_ThenProfileRefetchedForNewToken', async () => {
    usersApi.getMe
      .mockResolvedValueOnce({ firstName: 'Alice' })
      .mockResolvedValueOnce({ firstName: 'River' });

    localStorage.setItem('token', 'token-1');
    await useCurrentUser().load();

    expect(usersApi.getMe).toHaveBeenCalledTimes(1);
    expect(useCurrentUser().firstName.value).toBe('Alice');

    localStorage.setItem('token', 'token-2');
    await useCurrentUser().load();

    expect(usersApi.getMe).toHaveBeenCalledTimes(2);
    expect(useCurrentUser().firstName.value).toBe('River');
  });

  it('TestCurrentUser_GivenNoToken_WhenLoadCalled_ThenSkipsFetchAndClearsName', async () => {
    useCurrentUser().firstName.value = 'Alice';

    await useCurrentUser().load();

    expect(usersApi.getMe).not.toHaveBeenCalled();
    expect(useCurrentUser().firstName.value).toBe('');
  });
});
