import { renderWithProviders, screen } from '@/test/test-utils';
import { useDashboardActions, useDashboardPage } from '@/modules/dashboard/hooks';
import type { DashboardPageData } from '@/modules/dashboard/types';

import { DashboardPage } from '../DashboardPage';

vi.mock('@/modules/dashboard/hooks');

const mockUseDashboardPage = vi.mocked(useDashboardPage);
const mockUseDashboardActions = vi.mocked(useDashboardActions);

const createDashboardData = (overrides: Partial<DashboardPageData> = {}): DashboardPageData => ({
  activeRoom: null,
  activeRoomError: null,
  ledger: {
    achievements: [],
    currentLevelXp: 24,
    level: 2,
    nextLevelXp: 100,
    sessionsAdmined: 1,
    sessionsParticipated: 3,
    tasksEstimated: 8,
    xpProgressPercentage: 24
  },
  ledgerError: null,
  recentRooms: [],
  sessions: [],
  teams: [],
  teamsError: null,
  view: 'empty',
  ...overrides
});

describe('DashboardPage', () => {
  beforeEach(() => {
    mockUseDashboardActions.mockReturnValue({
      openCreateRoom: vi.fn(),
      openJoinRoom: vi.fn()
    });
  });

  it('renders a loading state while dashboard data is pending', () => {
    mockUseDashboardPage.mockReturnValue({
      data: null,
      errorMessage: null,
      retry: vi.fn(),
      status: 'loading'
    });

    renderWithProviders(<DashboardPage />);

    expect(screen.getByText('Loading dashboard')).toBeInTheDocument();
  });

  it('renders a retry state when the dashboard request fails', () => {
    const retry = vi.fn();

    mockUseDashboardPage.mockReturnValue({
      data: null,
      errorMessage: 'Backend unavailable.',
      retry,
      status: 'error'
    });

    renderWithProviders(<DashboardPage />);

    expect(screen.getByText('Dashboard unavailable')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Retry' })).toBeInTheDocument();
  });

  it('renders the empty dashboard composition', () => {
    const openCreateRoom = vi.fn();

    mockUseDashboardActions.mockReturnValue({
      openCreateRoom,
      openJoinRoom: vi.fn()
    });
    mockUseDashboardPage.mockReturnValue({
      data: createDashboardData(),
      errorMessage: null,
      retry: vi.fn(),
      status: 'ready'
    });

    renderWithProviders(<DashboardPage />);

    expect(screen.getByText('No active rooms yet')).toBeInTheDocument();
    expect(screen.getByText('No teams yet')).toBeInTheDocument();
  });

  it('renders the no-active-room composition with recent rooms only', () => {
    mockUseDashboardPage.mockReturnValue({
      data: createDashboardData({
        recentRooms: [
          {
            approxDurationSeconds: 1200,
            createdAt: '2026-04-05T10:00:00.000Z',
            estimatedTasksCount: 2,
            finishedAt: null,
            id: 'room-1',
            lastActivityAt: '2026-04-05T12:00:00.000Z',
            name: 'Core Auth Refactor',
            participantsCount: 5,
            role: 'ADMIN',
            status: 'ACTIVE',
            tasksCount: 3,
            teamId: null
          }
        ],
        sessions: [
          {
            approxDurationSeconds: 1200,
            createdAt: '2026-04-05T10:00:00.000Z',
            estimatedTasksCount: 2,
            finishedAt: null,
            id: 'room-1',
            lastActivityAt: '2026-04-05T12:00:00.000Z',
            name: 'Core Auth Refactor',
            participantsCount: 5,
            role: 'ADMIN',
            status: 'ACTIVE',
            tasksCount: 3,
            teamId: null
          }
        ],
        teams: [
          {
            createdAt: '2026-03-30T10:00:00.000Z',
            id: 'team-1',
            name: 'Platform Engineering',
            ownerUserId: 'user-1'
          }
        ],
        view: 'noActive'
      }),
      errorMessage: null,
      retry: vi.fn(),
      status: 'ready'
    });

    renderWithProviders(<DashboardPage />);

    expect(screen.getByText('Recent Rooms')).toBeInTheDocument();
    expect(screen.getByText('Core Auth Refactor')).toBeInTheDocument();
    expect(screen.getByText('Platform Engineering')).toBeInTheDocument();
  });

  it('renders the active room composition when an active room is available', () => {
    mockUseDashboardPage.mockReturnValue({
      data: createDashboardData({
        activeRoom: {
          code: '01HRXGQW4QK8M3C1A4Q0R2D8TM',
          currentTaskStatus: 'VOTING',
          currentTaskTitle: 'Kafka Event Bus Implementation',
          estimatedTasksCount: 3,
          id: 'room-3',
          lastActivityAt: '2026-04-06T14:00:00.000Z',
          name: 'Refinement: Microservices Core',
          participants: [
            {
              avatarUrl: null,
              displayName: 'Alex Architect',
              id: 'participant-1',
              role: 'ADMIN'
            }
          ],
          status: 'ACTIVE',
          tasksCount: 5,
          teamId: null
        },
        recentRooms: [
          {
            approxDurationSeconds: 1800,
            createdAt: '2026-04-03T10:00:00.000Z',
            estimatedTasksCount: 2,
            finishedAt: '2026-04-03T10:30:00.000Z',
            id: 'room-4',
            lastActivityAt: '2026-04-03T10:30:00.000Z',
            name: 'API Gateway Logic',
            participantsCount: 3,
            role: 'PARTICIPANT',
            status: 'FINISHED',
            tasksCount: 2,
            teamId: null
          }
        ],
        view: 'active'
      }),
      errorMessage: null,
      retry: vi.fn(),
      status: 'ready'
    });

    renderWithProviders(<DashboardPage />);

    expect(screen.getByText('Active Session')).toBeInTheDocument();
    expect(screen.getByText('Refinement: Microservices Core')).toBeInTheDocument();
    expect(screen.getByText('Recent Rooms')).toBeInTheDocument();
    expect(screen.getByText('API Gateway Logic')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Enter Room' })).toBeInTheDocument();
  });
});
