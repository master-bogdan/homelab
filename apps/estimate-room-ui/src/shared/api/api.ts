import { createApi } from '@reduxjs/toolkit/query/react';

import { baseQueryWithAuth } from './baseQuery';

export const api = createApi({
  baseQuery: baseQueryWithAuth,
  endpoints: () => ({}),
  reducerPath: 'api'
});
