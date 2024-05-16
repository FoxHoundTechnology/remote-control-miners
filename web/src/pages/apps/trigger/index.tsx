import { useEffect } from 'react'

import Grid from '@mui/material/Grid'
import Typography from '@mui/material/Typography'

import PageHeader from 'src/@core/components/page-header'

import TriggerCards from 'src/views/apps/trigger/TriggerCards'
import TimelineFilled from 'src/views/components/timeline/TimelineFilled'
import TriggerTimeline from 'src/views/apps/trigger/TriggerTimeline'
import { Trigger, TriggerHistory } from 'src/types/apps/triggerTypes'

// Redux related
import { useSelector, useDispatch } from 'react-redux'
import {
  FetchTriggers,
  CreateTrigger,
  UpdateTrigger,
  DeleteTrigger,
  RestartTrigger,
  StopTrigger
} from 'src/store/apps/trigger'
import { AppDispatch, RootState } from 'src/store'
import { Divider } from '@mui/material'

export interface ITriggerHandler {
  handleCreate: (trigger: Trigger) => void
  handleUpdate: (trigger: Trigger) => void
  handleDelete: (id: number) => void
  handleRestart: (trigger: Trigger) => void
  handleStop: (id: number) => void
}

const TriggerComponent = () => {
  // card selector

  // === Trigger Store Logic (Redux) ===
  const dispatch = useDispatch<AppDispatch>()
  const { data, status, message } = useSelector((state: RootState) => state.trigger)

  useEffect(() => {
    dispatch(FetchTriggers())
  }, [])

  const TriggerHandler: ITriggerHandler = {
    handleCreate: (trigger: Trigger) => {
      dispatch(CreateTrigger(trigger))
    },
    handleUpdate: (trigger: Trigger) => {
      dispatch(UpdateTrigger(trigger))
    },
    handleDelete: (id: number) => {
      dispatch(DeleteTrigger(id))
    },
    handleRestart: (trigger: Trigger) => {
      dispatch(RestartTrigger(trigger))
    },
    handleStop: (id: number) => {
      dispatch(StopTrigger(id))
    }
  }

  // and sync with the timeline
  return (
    <Grid container spacing={6}>
      <Grid item xs={12}>
        <Divider />
      </Grid>
      <Grid item xs={12} sx={{ mb: 4 }}>
        <TriggerCards triggers={data} triggerHandler={TriggerHandler} />
      </Grid>
      <Grid item xs={12}>
        <Divider />
      </Grid>
      <Grid item xs={12}>
        <TriggerTimeline triggers={data} />
        {/* <Table /> */}
      </Grid>
    </Grid>
  )
}

export default TriggerComponent
