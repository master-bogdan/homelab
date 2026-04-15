import { createTestStore } from '@/test/test-utils';

import { fetchDashboardPage, submitCreateRoom, submitJoinRoom } from '../dashboardThunks';

const createRoomDto = {
  adminUserId: 'user-1',
  code: '01HRXGQW4QK8M3C1A4Q0R2D8TM',
  createdAt: '2026-04-07T10:00:00.000Z',
  deck: {
    kind: 'FIBONACCI',
    name: 'Fibonacci',
    values: ['0', '1', '2', '3', '5', '8', '13', '21', '?']
  },
  finishedAt: null,
  lastActivityAt: '2026-04-07T10:05:00.000Z',
  name: 'Refinement: Microservices Core',
  participants: [
    {
      guestName: null,
      joinedAt: '2026-04-07T10:00:00.000Z',
      leftAt: null,
      role: 'ADMIN',
      roomId: 'room-1',
      roomParticipantId: 'participant-1',
      user: {
        avatarUrl: null,
        displayName: 'Alex Architect',
        email: 'alex@example.com',
        userId: 'user-1'
      },
      userId: 'user-1'
    }
  ],
  roomId: 'room-1',
  status: 'ACTIVE',
  tasks: [
    {
      createdAt: '2026-04-07T10:00:00.000Z',
      description: null,
      externalKey: null,
      finalEstimateValue: null,
      isActive: true,
      roomId: 'room-1',
      status: 'VOTING',
      taskId: 'task-1',
      title: 'Kafka Event Bus Implementation',
      updatedAt: '2026-04-07T10:05:00.000Z'
    },
    {
      createdAt: '2026-04-07T10:02:00.000Z',
      description: null,
      externalKey: null,
      finalEstimateValue: '5',
      isActive: false,
      roomId: 'room-1',
      status: 'ESTIMATED',
      taskId: 'task-2',
      title: 'Retry Strategy',
      updatedAt: '2026-04-07T10:05:00.000Z'
    }
  ],
  teamId: null
} as const;

const createJsonResponse = (payload: unknown, status = 200) =>
  new Response(JSON.stringify(payload), {
    headers: {
      'content-type': 'application/json'
    },
    status
  });

describe('dashboardThunks', () => {
  let fetchMock: ReturnType<typeof vi.fn>;

  beforeEach(() => {
    fetchMock = vi.fn();
    window.localStorage.clear();
    vi.stubGlobal('fetch', fetchMock);
  });

  afterEach(() => {
    vi.restoreAllMocks();
    vi.unstubAllGlobals();
  });

  it('maps dashboard page data into the active-room composition', async () => {
    fetchMock
      .mockResolvedValueOnce(createJsonResponse({
        items: [
          {
            approxDurationSeconds: 900,
            createdAt: '2026-04-07T10:00:00.000Z',
            estimatedTasksCount: 1,
            finishedAt: null,
            lastActivityAt: '2026-04-07T10:05:00.000Z',
            name: 'Refinement: Microservices Core',
            participantsCount: 4,
            role: 'ADMIN',
            roomId: 'room-1',
            status: 'ACTIVE',
            tasksCount: 2,
            teamId: null
          },
          {
            approxDurationSeconds: 1800,
            createdAt: '2026-04-04T10:00:00.000Z',
            estimatedTasksCount: 6,
            finishedAt: '2026-04-04T10:30:00.000Z',
            lastActivityAt: '2026-04-04T10:30:00.000Z',
            name: 'Schema Migration',
            participantsCount: 3,
            role: 'PARTICIPANT',
            roomId: 'room-2',
            status: 'FINISHED',
            tasksCount: 8,
            teamId: 'team-1'
          }
        ],
        page: 1,
        pageSize: 20,
        total: 2
      }))
      .mockResolvedValueOnce(createJsonResponse([
        {
          createdAt: '2026-03-30T10:00:00.000Z',
          name: 'Platform Engineering',
          ownerUserId: 'user-1',
          teamId: 'team-1'
        }
      ]))
      .mockResolvedValueOnce(createJsonResponse({
        achievements: [
          {
            key: 'SESSION_PARTICIPATION',
            level: 1,
            unlockedAt: '2026-04-01T10:00:00.000Z'
          }
        ],
        stats: {
          level: 2,
          nextLevelXp: 100,
          sessionsAdmined: 1,
          sessionsParticipated: 3,
          tasksEstimated: 8,
          xp: 24
        }
      }))
      .mockResolvedValueOnce(createJsonResponse(createRoomDto));

    const store = createTestStore();

    const result = await store.dispatch(fetchDashboardPage()).unwrap();

    expect(result.view).toBe('active');
    expect(result.activeRoom?.id).toBe('room-1');
    expect(result.activeRoom?.participants[0]?.displayName).toBe('Alex Architect');
    expect(result.teams[0]?.name).toBe('Platform Engineering');
    expect(result.ledger?.xpProgressPercentage).toBe(24);
    expect(result.recentRooms).toHaveLength(2);
  });

  it('preserves section-level errors when secondary requests fail', async () => {
    fetchMock
      .mockResolvedValueOnce(createJsonResponse({
        items: [
          {
            approxDurationSeconds: 900,
            createdAt: '2026-04-07T10:00:00.000Z',
            estimatedTasksCount: 1,
            finishedAt: null,
            lastActivityAt: '2026-04-07T10:05:00.000Z',
            name: 'Refinement: Microservices Core',
            participantsCount: 4,
            role: 'ADMIN',
            roomId: 'room-1',
            status: 'ACTIVE',
            tasksCount: 2,
            teamId: null
          }
        ],
        page: 1,
        pageSize: 20,
        total: 1
      }))
      .mockResolvedValueOnce(
        createJsonResponse({ detail: 'teams down', status: 500, title: 'Internal Error' }, 500)
      )
      .mockResolvedValueOnce(
        createJsonResponse({ detail: 'ledger down', status: 500, title: 'Internal Error' }, 500)
      )
      .mockResolvedValueOnce(
        createJsonResponse({ detail: 'room down', status: 500, title: 'Internal Error' }, 500)
      );

    const store = createTestStore();

    const result = await store.dispatch(fetchDashboardPage()).unwrap();

    expect(result.view).toBe('noActive');
    expect(result.activeRoom).toBeNull();
    expect(result.activeRoomError?.message).toContain('room down');
    expect(result.teamsError?.message).toContain('teams down');
    expect(result.ledgerError?.message).toContain('ledger down');
  });

  it('maps the create-room response into the success dialog result', async () => {
    fetchMock.mockResolvedValueOnce(createJsonResponse({
      inviteToken: 'fallback-token-1',
      room: createRoomDto,
      shareLink: {
        acceptedAt: null,
        createdAt: '2026-04-07T10:00:00.000Z',
        createdByUserId: 'user-1',
        declinedAt: null,
        invitedEmail: null,
        invitedUserId: null,
        invitationId: 'invite-1',
        kind: 'ROOM_LINK',
        revokedAt: null,
        roomId: 'room-1',
        status: 'ACTIVE',
        teamId: null,
        token: 'share-token-1',
        updatedAt: '2026-04-07T10:00:00.000Z'
      },
      skippedRecipients: [
        {
          email: 'already@joined.dev',
          reason: 'already_joined',
          userId: null
        }
      ]
    }));

    const store = createTestStore();

    const result = await store.dispatch(
      submitCreateRoom({
        createShareLink: true,
        deckKey: 'fibonacci',
        inviteEmails: 'first@example.com, second@example.com',
        inviteTeamId: '',
        name: 'Refinement: Microservices Core'
      })
    ).unwrap();

    expect(result.roomCode).toBe('share-token-1');
    expect(result.inviteLink).toBe('http://localhost:3000/join/share-token-1');
    expect(result.skippedRecipients).toHaveLength(1);
  });

  it('falls back to inviteToken when no share link object is returned', async () => {
    fetchMock.mockResolvedValueOnce(createJsonResponse({
      inviteToken: 'fallback-token-1',
      room: createRoomDto
    }));

    const store = createTestStore();

    const result = await store.dispatch(
      submitCreateRoom({
        createShareLink: true,
        deckKey: 'fibonacci',
        inviteEmails: '',
        inviteTeamId: '',
        name: 'Refinement: Microservices Core'
      })
    ).unwrap();

    expect(result.roomCode).toBe('fallback-token-1');
    expect(result.inviteLink).toBe('http://localhost:3000/join/fallback-token-1');
  });

  it('joins a room from a pasted invite URL', async () => {
    fetchMock.mockResolvedValueOnce(createJsonResponse({
      acceptedAt: null,
      createdAt: '2026-04-07T10:00:00.000Z',
      createdByUserId: 'user-1',
      declinedAt: null,
      invitedEmail: null,
      invitedUserId: null,
      invitationId: 'invite-1',
      kind: 'ROOM_LINK',
      revokedAt: null,
      roomId: 'room-1',
      status: 'ACTIVE',
      teamId: null,
      updatedAt: '2026-04-07T10:00:00.000Z'
    }));
    fetchMock.mockResolvedValueOnce(createJsonResponse({
      participant: createRoomDto.participants[0],
      room: createRoomDto
    }));

    const store = createTestStore();

    const result = await store.dispatch(
      submitJoinRoom('https://app.example.com/invites/token-123')
    ).unwrap();

    expect(result.roomId).toBe('room-1');
    expect(result.roomName).toBe('Refinement: Microservices Core');
  });

  it('rejects team invitations in the dashboard join-room flow', async () => {
    fetchMock.mockResolvedValueOnce(createJsonResponse({
      acceptedAt: null,
      createdAt: '2026-04-07T10:00:00.000Z',
      createdByUserId: 'user-1',
      declinedAt: null,
      invitedEmail: 'member@example.com',
      invitedUserId: null,
      invitationId: 'invite-2',
      kind: 'TEAM_MEMBER',
      revokedAt: null,
      roomId: null,
      status: 'ACTIVE',
      teamId: 'team-1',
      updatedAt: '2026-04-07T10:00:00.000Z'
    }));

    const store = createTestStore();

    await expect(store.dispatch(submitJoinRoom('team-code-1')).unwrap()).rejects.toBe(
      'This code belongs to a team invitation, not a room session.'
    );
  });
});
