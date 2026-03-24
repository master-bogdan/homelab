export const appRoutes = {
  root: '/',
  login: '/login',
  dashboard: '/dashboard',
  roomsNew: '/rooms/new',
  roomDetails: '/rooms/:id',
  roomDetailsPath: (id: string) => `/rooms/${id}`,
  history: '/history',
  historyRoom: '/history/rooms/:id',
  historyRoomPath: (id: string) => `/history/rooms/${id}`,
  teamDetails: '/teams/:id',
  teamDetailsPath: (id: string) => `/teams/${id}`,
  profile: '/profile',
  settings: '/settings'
} as const;
