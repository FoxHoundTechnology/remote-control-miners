import React from 'react'
import { Grid } from '@mui/material'

import { ActionType, TargetType, Trigger } from 'src/types/apps/triggerTypes'
import TriggerCard from './TriggerCard'
import { ITriggerHandler } from 'src/pages/apps/trigger'

interface TriggerCardsProps {
  triggers: Trigger[]
  triggerHandler: ITriggerHandler
}

const initialStateForCurrentTrigger = {
  ID: 0,
  name: '',
  user_id: 0, // FIXME after user auth is fully implemented
  interval: 10,
  active: false,
  last_executed: new Date(),
  targets: [
    {
      type: TargetType.Hashrate,
      percentage: 0,
      value: 0
    }
  ],
  actions: [
    {
      type: ActionType.Reboot,
      interval: 0,
      value: ''
    }
  ]
}

const TriggerCards: React.FC<TriggerCardsProps> = ({ triggers, triggerHandler }) => {
  return (
    <>
      <Grid container spacing={3}>
        {triggers.map((trigger, index) => (
          <Grid item xs={12} sm={6} md={4} key={index}>
            <TriggerCard trigger={trigger} triggerHandler={triggerHandler} />
          </Grid>
        ))}

        {/* Render the "Add New Trigger" card */}
        <Grid item xs={12} sm={6} md={4}>
          <TriggerCard trigger={initialStateForCurrentTrigger} addNew triggerHandler={triggerHandler} />
        </Grid>
      </Grid>
    </>
  )
}

export default TriggerCards
