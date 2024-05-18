import { createSlice, createAsyncThunk } from '@reduxjs/toolkit'
import axios from 'axios'

import dotenv from 'dotenv'

import { Trigger } from 'src/types/apps/triggerTypes'

dotenv.config()

const apiServiceUrl = process.env.REMOTE_CONTROL_SERVER_URL

// Create an Axios instance for the service
const apiService = axios.create({
  baseURL: `${apiServiceUrl}/trigger`
})

// export interface TriggerState {
//   data: Trigger[]
//   status: 'N/A' | 'pending' | 'succeeded' | 'failed' // Add other statuses if needed.
//   message: string
// }

export const FetchTriggers = createAsyncThunk('triggers/FetchTrigers', async () => {
  const response = await apiService.get('/list')

  return response.data
})

export const CreateTrigger = createAsyncThunk('triggers/CreateTrigger', async (trigger: Trigger) => {
  const triggerWithoutID = {
    name: trigger.name,
    interval: trigger.interval,
    active: trigger.active,
    targets: trigger.targets,
    actions: trigger.actions,
    last_executed: trigger.last_executed,
    user_id: trigger.user_id
  }

  const response = await apiService.post('', triggerWithoutID)

  return response.data
})

export const UpdateTrigger = createAsyncThunk('triggers/UpdateTrigger', async (trigger: Trigger) => {
  const response = await apiService.put(``, trigger)

  return response.data
})

export const DeleteTrigger = createAsyncThunk('triggers/DeleteTrigger', async (triggerId: number) => {
  const response = await apiService
    .delete(``, {
      // NOTE: in delete request, the data is passed in the body with data key
      data: {
        id: triggerId
      }
    })
    .catch(err => {
      console.log('error from DeleteTrigger in catch', err)
    })

  console.log('response from DeleteTrigger', response)

  // this id will be passed to the reducer and the reducer will remove the trigger with this id
  return triggerId
})

export const RestartTrigger = createAsyncThunk('triggers/RestartTrigger', async (trigger: Trigger) => {
  const response = await apiService.post(`/start`, trigger)
  console.log('response from RestartTrigger in createAsync', response)

  return trigger.ID
})

export const StopTrigger = createAsyncThunk('triggers/StopTrigger', async (triggerId: number) => {
  console.log('passed id is', triggerId)

  await apiService
    .post(`/stop`, {
      id: triggerId
    })
    .then(res => {
      console.log('res from StopTrigger in createAsync', res)
    })

  return triggerId
})

export const DeleteTriggerHistories = createAsyncThunk('triggers/DeleteTriggerHistories', async (triggerId: number) => {
  console.log('passed id is', triggerId)

  await apiService.delete(`/history`, {
    data: {
      id: triggerId
    }
  })

  return triggerId
})

export const TriggersSlice = createSlice({
  name: 'triggers',
  initialState: {
    data: [] as Trigger[],
    status: 'N/A',
    message: 'N/A'
  },
  reducers: {},
  extraReducers: builder => {
    // Fetch Triggers
    builder.addCase(FetchTriggers.fulfilled, (state, action) => {
      state.status = 'succeeded'
      state.data = action.payload?.data
      console.log('action?.payload.data from FetchTriggers', action.payload?.data)
    })
    builder.addCase(FetchTriggers.pending, state => {
      state.status = 'pending'
    })
    builder.addCase(FetchTriggers.rejected, state => {
      state.status = 'failed'
    })

    // Create Trigger
    builder.addCase(CreateTrigger.fulfilled, (state, action) => {
      state.status = 'succeeded'

      const newTrigger = action.payload?.data
      console.log('newTrigger from CreateTrigger builder', newTrigger)
      // assuming that the response will be the newly created trigger
      state.data.push(newTrigger)
    })
    builder.addCase(CreateTrigger.pending, state => {
      state.status = 'pending'
    })
    builder.addCase(CreateTrigger.rejected, state => {
      state.status = 'failed'
    })

    // Update Trigger
    builder.addCase(UpdateTrigger.fulfilled, (state, action) => {
      console.log('success action.payload?.data from UpdateTrigger', action.payload?.data)
      state.status = 'succeeded'
      state.data = state.data.map((trigger: Trigger) => {
        return trigger.ID === action.payload?.data.ID ? action.payload?.data : trigger
      })
    })
    builder.addCase(UpdateTrigger.pending, state => {
      state.status = 'pending'
    })
    builder.addCase(UpdateTrigger.rejected, state => {
      state.status = 'failed'
    })

    // Delete Trigger
    builder.addCase(DeleteTrigger.fulfilled, (state, action) => {
      state.status = 'succeeded'
      if (Array.isArray(state.data)) {
        state.data = state.data.filter((trigger: Trigger) => trigger.ID !== action.payload)
      } else {
        console.log('state.data is not an array:', state.data)
      }
    })

    builder.addCase(DeleteTrigger.pending, state => {
      state.status = 'pending'
    })

    builder.addCase(DeleteTrigger.rejected, state => {
      state.status = 'failed'
    })

    // Restart Trigger
    builder.addCase(RestartTrigger.fulfilled, (state, action) => {
      state.status = 'succeeded'
      console.log('action.payload from RestartTrigger', action.payload)
      state.data = state.data.map((trigger: Trigger) => {
        console.log('trigger.ID from RestartTrigger', trigger.ID)
        console.log('action.payload from RestartTrigger', action.payload)
        return trigger.ID === action.payload
          ? {
              ...trigger,
              active: true
            }
          : trigger
      })
    })

    builder.addCase(RestartTrigger.pending, state => {
      state.status = 'pending'
    })

    builder.addCase(RestartTrigger.rejected, state => {
      state.status = 'failed'
    })

    // Stop Trigger
    builder.addCase(StopTrigger.fulfilled, (state, action) => {
      state.status = 'succeeded'
      state.data = state.data.map((trigger: Trigger) => {
        return trigger.ID === action.payload
          ? {
              ...trigger,
              active: false
            }
          : trigger
      })
    })

    builder.addCase(StopTrigger.pending, state => {
      state.status = 'pending'
    })

    builder.addCase(StopTrigger.rejected, state => {
      state.status = 'failed'
    })

    // Delte Trigger Histories
    builder.addCase(DeleteTriggerHistories.fulfilled, (state, action) => {
      state.status = 'succeeded'
      state.data = state.data.map((trigger: Trigger) => {
        return trigger.ID === action.payload
          ? {
              ...trigger,
              histories: []
            }
          : trigger
      })
    })

    builder.addCase(DeleteTriggerHistories.pending, state => {
      state.status = 'pending'
    })

    builder.addCase(DeleteTriggerHistories.rejected, state => {
      state.status = 'failed'
    })
  }
})

export default TriggersSlice.reducer
