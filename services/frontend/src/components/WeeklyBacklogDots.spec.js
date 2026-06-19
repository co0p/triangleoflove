import { describe, it, expect } from 'vitest';
import { mount } from '@vue/test-utils';
import WeeklyBacklogDots from './WeeklyBacklogDots.vue';

describe('WeeklyBacklogDots', () => {
  it('renders 3 dot slots', () => {
    const wrapper = mount(WeeklyBacklogDots, {
      props: { weeklyInsights: [] },
    });
    const dots = wrapper.findAll('[data-testid="weekly-backlog-dot"]');
    expect(dots).toHaveLength(3);
  });

  it('shows 3 dots when no insights provided', () => {
    const wrapper = mount(WeeklyBacklogDots, {
      props: { weeklyInsights: [] },
    });
    const filledDots = wrapper.findAll('.weekly-backlog-dot--filled');
    expect(filledDots).toHaveLength(3);
  });

  it('shows 1 dot for today entry', () => {
    const now = new Date('2026-06-19T12:00:00Z');
    const wrapper = mount(WeeklyBacklogDots, {
      props: {
        weeklyInsights: [{ date: '20260619', value: 10 }],
        referenceDate: now,
      },
    });
    const filledDots = wrapper.findAll('.weekly-backlog-dot--filled');
    expect(filledDots).toHaveLength(1);
  });

  it('shows 1 dot for yesterday entry', () => {
    const now = new Date('2026-06-19T12:00:00Z');
    const wrapper = mount(WeeklyBacklogDots, {
      props: {
        weeklyInsights: [{ date: '20260618', value: 10 }],
        referenceDate: now,
      },
    });
    const filledDots = wrapper.findAll('.weekly-backlog-dot--filled');
    expect(filledDots).toHaveLength(1);
  });

  it('shows 2 dots for 2 days ago', () => {
    const now = new Date('2026-06-19T12:00:00Z');
    const wrapper = mount(WeeklyBacklogDots, {
      props: {
        weeklyInsights: [{ date: '20260617', value: 10 }],
        referenceDate: now,
      },
    });
    const filledDots = wrapper.findAll('.weekly-backlog-dot--filled');
    expect(filledDots).toHaveLength(2);
  });

  it('shows 2 dots for 3 days ago', () => {
    const now = new Date('2026-06-19T12:00:00Z');
    const wrapper = mount(WeeklyBacklogDots, {
      props: {
        weeklyInsights: [{ date: '20260616', value: 10 }],
        referenceDate: now,
      },
    });
    const filledDots = wrapper.findAll('.weekly-backlog-dot--filled');
    expect(filledDots).toHaveLength(2);
  });

  it('shows 3 dots for 4 days ago', () => {
    const now = new Date('2026-06-19T12:00:00Z');
    const wrapper = mount(WeeklyBacklogDots, {
      props: {
        weeklyInsights: [{ date: '20260615', value: 10 }],
        referenceDate: now,
      },
    });
    const filledDots = wrapper.findAll('.weekly-backlog-dot--filled');
    expect(filledDots).toHaveLength(3);
  });

  it('uses latest date when multiple entries exist', () => {
    const now = new Date('2026-06-19T12:00:00Z');
    const wrapper = mount(WeeklyBacklogDots, {
      props: {
        weeklyInsights: [
          { date: '20260615', value: 10 },
          { date: '20260617', value: 8 },
          { date: '20260612', value: 5 },
        ],
        referenceDate: now,
      },
    });
    const filledDots = wrapper.findAll('.weekly-backlog-dot--filled');
    expect(filledDots).toHaveLength(2); // 2 days behind (from 20260617)
  });

  it('has correct aria-label for 1 dot', () => {
    const now = new Date('2026-06-19T12:00:00Z');
    const wrapper = mount(WeeklyBacklogDots, {
      props: {
        weeklyInsights: [{ date: '20260619', value: 10 }],
        referenceDate: now,
      },
    });
    expect(wrapper.attributes('aria-label')).toBe('Weekly check-in backlog: 1 of 3 dots');
  });

  it('has correct aria-label for 2 dots', () => {
    const now = new Date('2026-06-19T12:00:00Z');
    const wrapper = mount(WeeklyBacklogDots, {
      props: {
        weeklyInsights: [{ date: '20260617', value: 10 }],
        referenceDate: now,
      },
    });
    expect(wrapper.attributes('aria-label')).toBe('Weekly check-in backlog: 2 of 3 dots');
  });

  it('has correct aria-label for 3 dots', () => {
    const wrapper = mount(WeeklyBacklogDots, {
      props: { weeklyInsights: [] },
    });
    expect(wrapper.attributes('aria-label')).toBe('Weekly check-in backlog: 3 of 3 dots');
  });

  it('has correct testid on root', () => {
    const wrapper = mount(WeeklyBacklogDots, {
      props: { weeklyInsights: [] },
    });
    expect(wrapper.attributes('data-testid')).toBe('weekly-backlog-dots');
  });

  it('has role="img" for semantic purpose', () => {
    const wrapper = mount(WeeklyBacklogDots, {
      props: { weeklyInsights: [] },
    });
    expect(wrapper.attributes('role')).toBe('img');
  });

  it('marks individual dots as aria-hidden', () => {
    const wrapper = mount(WeeklyBacklogDots, {
      props: { weeklyInsights: [] },
    });
    const dots = wrapper.findAll('[data-testid="weekly-backlog-dot"]');
    dots.forEach(dot => {
      expect(dot.attributes('aria-hidden')).toBe('true');
    });
  });
});
