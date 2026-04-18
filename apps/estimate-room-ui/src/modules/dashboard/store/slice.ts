import { createSlice } from '@reduxjs/toolkit';

import type {
  DashboardCreateRoomState,
  DashboardJoinRoomState,
  DashboardPageState,
  DashboardState
} from '../types';
import {
  fetchCreateRoomTeams,
  fetchDashboardPage,
  submitCreateRoom,
  submitJoinRoom
} from './thunks';

const initialPageState: DashboardPageState = {
  data: null,
  errorMessage: null,
  status: 'loading'
};

const initialCreateRoomState: DashboardCreateRoomState = {
  isLoadingTeams: false,
  submitErrorMessage: null,
  teamErrorMessage: null,
  teamOptions: []
};

const initialJoinRoomState: DashboardJoinRoomState = {
  errorMessage: null
};

const initialState: DashboardState = {
  createRoom: initialCreateRoomState,
  joinRoom: initialJoinRoomState,
  page: initialPageState
};

const dashboardSlice = createSlice({
  name: 'dashboard',
  initialState,
  reducers: {
    resetCreateRoomDialogState: (state) => {
      state.createRoom.submitErrorMessage = null;
      state.createRoom.teamErrorMessage = null;
    },
    resetJoinRoomDialogState: (state) => {
      state.joinRoom.errorMessage = null;
    }
  },
  extraReducers: (builder) => {
    builder
      .addCase(fetchDashboardPage.pending, (state) => {
        state.page.errorMessage = null;
        state.page.status = 'loading';
      })
      .addCase(fetchDashboardPage.fulfilled, (state, action) => {
        state.page.data = action.payload;
        state.page.errorMessage = null;
        state.page.status = 'ready';
      })
      .addCase(fetchDashboardPage.rejected, (state, action) => {
        state.page.data = null;
        state.page.errorMessage =
          typeof action.payload === 'string'
            ? action.payload
            : 'Dashboard data could not be loaded right now.';
        state.page.status = 'error';
      })
      .addCase(fetchCreateRoomTeams.pending, (state) => {
        state.createRoom.isLoadingTeams = true;
        state.createRoom.teamErrorMessage = null;
      })
      .addCase(fetchCreateRoomTeams.fulfilled, (state, action) => {
        state.createRoom.isLoadingTeams = false;
        state.createRoom.teamErrorMessage = null;
        state.createRoom.teamOptions = action.payload;
      })
      .addCase(fetchCreateRoomTeams.rejected, (state, action) => {
        state.createRoom.isLoadingTeams = false;
        state.createRoom.teamErrorMessage =
          typeof action.payload === 'string'
            ? action.payload
            : 'Teams could not be loaded for room creation.';
        state.createRoom.teamOptions = [];
      })
      .addCase(submitCreateRoom.pending, (state) => {
        state.createRoom.submitErrorMessage = null;
      })
      .addCase(submitCreateRoom.fulfilled, (state) => {
        state.createRoom.submitErrorMessage = null;
      })
      .addCase(submitCreateRoom.rejected, (state, action) => {
        state.createRoom.submitErrorMessage =
          typeof action.payload === 'string'
            ? action.payload
            : 'The room could not be created right now.';
      })
      .addCase(submitJoinRoom.pending, (state) => {
        state.joinRoom.errorMessage = null;
      })
      .addCase(submitJoinRoom.fulfilled, (state) => {
        state.joinRoom.errorMessage = null;
      })
      .addCase(submitJoinRoom.rejected, (state, action) => {
        state.joinRoom.errorMessage =
          typeof action.payload === 'string'
            ? action.payload
            : 'Invalid or expired room code. Please check and try again.';
      });
  }
});

export const { resetCreateRoomDialogState, resetJoinRoomDialogState } =
  dashboardSlice.actions;
export const dashboardReducer = dashboardSlice.reducer;
