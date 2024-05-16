import React, { useEffect, useState } from 'react'

import {
  Button,
  TextField,
  FormControlLabel,
  Switch,
  Container,
  Grid,
  Card,
  CardContent,
  Chip,
  Typography,
  useTheme,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  Divider,
  styled,
  Skeleton
} from '@mui/material'
import { Target, Action, Trigger, ActionType, TargetType } from 'src/types/apps/triggerTypes'

import { useFormik } from 'formik'
import * as Yup from 'yup'

import { ITriggerHandler } from 'src/pages/apps/trigger'
import ConfirmModal from './ConfirmModal'
import { useSelector } from 'react-redux'
import { RootState } from 'src/store'
import { generateTargetSentence, generateActionSentence, generateValueLabel } from './utils'

type TriggerCardProps = {
  trigger: Trigger
  addNew?: boolean // This is just to differentiate between a new trigger card and an existing trigger card
  triggerHandler: ITriggerHandler
}

const initialTarget = {
  type: TargetType.Hashrate,
  percentage: 0,
  value: 0
}

const initialAction = {
  type: ActionType.Reboot,
  interval: 0,
  value: ''
}

const cardMinHeight = '185px' // Adjust this value as needed

const AndDivider: React.FC = () => (
  <div style={{ display: 'flex', alignItems: 'center', margin: '0px 0' }}>
    <Divider style={{ flex: 1 }} />
    <span style={{ margin: '0px', color: 'lightgrey' }}>
      <Typography variant='caption'>AND</Typography>
    </span>
    <Divider style={{ flex: 1 }} />
  </div>
)

const StyledChip = styled(Chip)(({ theme }) => ({
  backgroundColor: theme.palette.primary.main,
  color: theme.palette.primary.contrastText,
  fontSize: '1.03em', // Adjust this as per your need
  fontWeight: 700
}))

// TODO: UI improvements
//       by utilizing the status from the redux store, we can show a loading skeleton
const skeltonLoading = () => {
  return (
    <Card>
      <CardContent>
        <Skeleton variant='text' width='60%' />
        <Skeleton variant='text' width='80%' height={40} />
        <Skeleton variant='text' width='50%' />
        <Skeleton variant='text' width='40%' />
      </CardContent>
    </Card>
  )
}

const TriggerCard = ({ trigger, addNew, triggerHandler }: TriggerCardProps) => {
  console.log('trigger id', trigger.ID)
  const formik = useFormik({
    initialValues: trigger,
    validationSchema: Yup.object({
      name: Yup.string().required('Required'),
      interval: Yup.number().required('Required').min(10, 'Interval should be at least 10')
      // Add more validation rules as needed
    }),
    onSubmit: values => {
      // formik.resetForm()
    }
  })

  const { status } = useSelector((state: RootState) => state.trigger)

  useEffect(() => {}, [status, trigger])

  const theme = useTheme()
  const [open, setOpen] = useState(false)
  const [confirmModal, setConfirmModal] = useState(false)

  const handleClickOpen = () => {
    setOpen(true)
  }

  const handleClose = () => {
    console.log('closed')
    formik.resetForm()

    setOpen(false)
  }

  const handleCreate = () => {
    triggerHandler.handleCreate(formik.values)
    setOpen(false)
  }

  const handleUpdate = () => {
    triggerHandler.handleUpdate(formik.values)
    setOpen(false)
  }

  const handleDelete = () => {
    if (trigger.ID !== undefined) {
      triggerHandler.handleDelete(trigger.ID)
      setOpen(false)
    }
  }

  return (
    <>
      {addNew === true ? (
        // for the "Add New Trigger" card
        <Card
          onClick={() => {
            formik.resetForm()
            handleClickOpen()
          }}
          style={{
            cursor: 'pointer',
            transition: '0.3s',
            transform: 'scale(1)',
            border: '2px dashed lightgray',
            height: '155px',
            minHeight: cardMinHeight,
            display: 'flex',
            justifyContent: 'center',
            alignItems: 'center'
          }}
          onMouseOver={e => (e.currentTarget.style.transform = 'scale(1.05)')}
          onMouseOut={e => (e.currentTarget.style.transform = 'scale(1)')}
        >
          <CardContent>
            <Typography variant='h6' align='center'>
              + Add New Trigger
            </Typography>
          </CardContent>
        </Card>
      ) : (
        // for the existing trigger cards
        <Card
          onClick={() => handleClickOpen()}
          style={{
            cursor: 'pointer',
            transition: '0.3s',
            transform: 'scale(1)',
            minHeight: cardMinHeight,
            border: trigger.active ? `2px solid ${theme.palette.primary.main}` : 'none'
          }}
          onMouseOver={e => (e.currentTarget.style.transform = 'scale(1.05)')}
          onMouseOut={e => (e.currentTarget.style.transform = 'scale(1)')}
        >
          <CardContent>
            <Chip
              label={trigger.active ? 'Active' : 'Inactive'}
              color={trigger.active ? 'primary' : 'default'}
              variant='outlined'
              size='small'
              style={{ marginBottom: '4px' }}
            />
            <Typography variant='h6' style={{ marginBottom: '4px' }}>
              {trigger.name}
            </Typography>
            <Typography variant='caption' style={{ marginBottom: '4px' }}>
              Last Executed: {new Date(trigger.last_executed).toLocaleString()}
            </Typography>
            <Typography variant='body2' style={{ marginBottom: '4px' }}>
              Interval: {trigger.interval} mins
            </Typography>
          </CardContent>
        </Card>
      )}
      <Dialog fullWidth maxWidth='md' onClose={handleClose} open={open}>
        <DialogTitle style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
          {!addNew ? <span>Edit {formik.values.name}</span> : <span>Add New Trigger</span>}
          {!addNew && (
            <div style={{ marginLeft: 'auto' }}>
              <Button
                variant='outlined'
                onClick={() => {
                  setConfirmModal(true)
                }}
              >
                {trigger.active ? 'Disable this trigger' : 'Enable this trigger'}
              </Button>
            </div>
          )}
        </DialogTitle>
        <DialogContent>
          <form onSubmit={formik.handleSubmit}>
            <Grid container spacing={4}>
              <Grid item xs={6} padding={4}>
                <TextField
                  fullWidth
                  id='name'
                  name='name'
                  label='Name'
                  variant='outlined'
                  onChange={formik.handleChange}
                  defaultValue={trigger.name}
                  value={formik.values.name}
                  error={formik.touched.name && Boolean(formik.errors.name)}
                  helperText={formik.touched.name && formik.errors.name}
                  style={{ marginTop: '4px' }}
                />
              </Grid>
              <Grid item xs={6}>
                <TextField
                  fullWidth
                  id='interval'
                  name='interval'
                  label='Interval'
                  variant='outlined'
                  type='number'
                  defaultValue={trigger.interval}
                  onChange={formik.handleChange}
                  value={formik.values.interval}
                  error={formik.touched.interval && Boolean(formik.errors.interval)}
                  helperText={formik.touched.interval && formik.errors.interval}
                  style={{ marginTop: '4px' }}
                />
              </Grid>
              <Grid item xs={12}>
                <Typography variant='h6'>Targets</Typography>
                <Typography variant='caption'>
                  Please ensure that the sentence appears for every target and action setting.
                </Typography>
              </Grid>
              {formik.values.targets.map((target, index) => (
                <>
                  {generateTargetSentence(target) !== '' && (
                    <Grid item xs={12}>
                      <StyledChip label={generateTargetSentence(target)} />
                    </Grid>
                  )}

                  <Grid item xs={6} key={index}>
                    <FormControl fullWidth variant='outlined' style={{ marginTop: '4px' }}>
                      <InputLabel htmlFor={`target-type-${index}`}>Type</InputLabel>
                      <Select
                        label='Type'
                        id={`target-type-${index}`}
                        name={`targets.${index}.type`}
                        value={target.type}
                        onChange={formik.handleChange}
                        defaultValue={trigger?.targets[index]?.type}
                      >
                        {Object.values(TargetType).map(typeValue => (
                          <MenuItem key={typeValue} value={typeValue}>
                            {typeValue}
                          </MenuItem>
                        ))}
                      </Select>
                      {/* {formik.touched.targets?.[index]?.type && Boolean(formik.errors.targets?.[index]?.type) && (
                        <Typography color='error' variant='caption'>
                          ERROR
                        </Typography>
                      )} */}
                    </FormControl>
                  </Grid>
                  <Grid item xs={6}>
                    <TextField
                      fullWidth
                      id={`target-percentage-${index}`}
                      name={`targets.${index}.percentage`}
                      label='Pecentage of fleet'
                      variant='outlined'
                      type='number'
                      onChange={formik.handleChange}
                      defaultValue={trigger?.targets[index]?.percentage}
                      value={target.percentage || 0}
                      // error={formik.touched.targets?.[index]?.value && Boolean(formik.errors.targets?.[index]?.value)}
                      // helperText={formik.touched.targets?.[index]?.value && formik.errors.targets?.[index]?.value}
                      style={{ marginTop: '4px' }}
                    />
                  </Grid>
                  {target.type !== TargetType.Offline && target.type !== TargetType.MissingHashboard && (
                    <Grid item xs={6}>
                      <TextField
                        fullWidth
                        id={`target-value-${index}`}
                        name={`targets.${index}.value`}
                        label={generateValueLabel(target)}
                        variant='outlined'
                        type='number'
                        onChange={formik.handleChange}
                        value={target.value || 0}
                        defaultValue={trigger?.targets[index]?.value}
                        // error={formik.touched.targets?.[index]?.value && Boolean(formik.errors.targets?.[index]?.value)}
                        // helperText={formik.touched.targets?.[index]?.value && formik.errors.targets?.[index]?.value}
                        style={{ marginTop: '4px' }}
                      />
                    </Grid>
                  )}
                  {index !== formik.values.targets.length - 1 && (
                    <Grid item xs={12}>
                      <AndDivider />
                    </Grid>
                  )}
                </>
              ))}
              <Grid item xs={12} style={{ display: 'flex', justifyContent: 'flex-end' }}>
                <Button
                  variant='outlined'
                  color='primary'
                  onClick={() => {
                    const newTargets = [...formik.values.targets]
                    newTargets.push(initialTarget)
                    formik.setFieldValue('targets', newTargets)
                  }}
                  style={{ marginRight: '8px' }}
                >
                  ADD TARGET
                </Button>
                <Button
                  variant='outlined'
                  color='secondary'
                  onClick={() => {
                    const newTargets = [...formik.values.targets]
                    if (newTargets.length > 1) {
                      newTargets.pop() // Remove the last target
                    }
                    formik.setFieldValue('targets', newTargets)
                  }}
                  disabled={formik.values.targets.length === 1} // Disable the button if only one target
                >
                  REMOVE TARGET
                </Button>
              </Grid>
              <Grid item xs={12}>
                <Typography variant='h6'>Actions</Typography>
              </Grid>
              {formik.values.actions.map((action, index) => (
                <>
                  {generateActionSentence(action) !== '' && (
                    <Grid item xs={12}>
                      <StyledChip label={generateActionSentence(action)} />
                    </Grid>
                  )}
                  {/* Action Type dropdown */}
                  <Grid item xs={6}>
                    <FormControl fullWidth variant='outlined' style={{ marginTop: '16px' }}>
                      <InputLabel htmlFor={`action-type-${index}`}>Action Type</InputLabel>
                      <Select
                        label='Action Type'
                        id={`action-type-${index}`}
                        name={`actions.${index}.type`}
                        value={action.type}
                        onChange={formik.handleChange}
                        defaultValue={trigger?.actions[index]?.type}
                      >
                        {Object.values(ActionType).map(actionTypeValue => (
                          <MenuItem key={actionTypeValue} value={actionTypeValue}>
                            {actionTypeValue}
                          </MenuItem>
                        ))}
                      </Select>
                    </FormControl>
                  </Grid>

                  {/* Conditionally render Action Interval only when there are multiple actions */}
                  {formik.values.actions.length > 1 && index !== formik.values.actions.length - 1 && (
                    <Grid item xs={6}>
                      <TextField
                        fullWidth
                        id={`action-interval-${index}`}
                        name={`actions.${index}.interval`}
                        label='Interval'
                        variant='outlined'
                        type='number'
                        onChange={formik.handleChange}
                        value={action.interval || ''}
                        defaultValue={trigger?.actions[index]?.interval}
                        style={{ marginTop: '16px' }}
                      />
                    </Grid>
                  )}
                  {/* Add the divider, but not for the last action */}
                  {index !== formik.values.actions.length - 1 && (
                    <Grid item xs={12}>
                      <AndDivider />
                    </Grid>
                  )}
                </>
              ))}
              <Grid item xs={12} style={{ display: 'flex', justifyContent: 'flex-end' }}>
                <Button
                  variant='outlined'
                  color='primary'
                  onClick={() => {
                    const newActions = [...formik.values.actions]
                    newActions.push(initialAction) // Assuming you have a default structure for a new action
                    formik.setFieldValue('actions', newActions)
                  }}
                  style={{ marginRight: '8px' }}
                >
                  ADD ACTION
                </Button>
                <Button
                  variant='outlined'
                  color='secondary'
                  onClick={() => {
                    const newActions = [...formik.values.actions]
                    if (newActions.length > 1) {
                      newActions.pop() // Remove the last action
                    }
                    formik.setFieldValue('actions', newActions)
                  }}
                  disabled={formik.values.actions.length === 1} // Disable the button if only one action
                >
                  REMOVE ACTION
                </Button>
              </Grid>
            </Grid>
          </form>
        </DialogContent>
        <DialogActions>
          {!trigger?.active && (
            <Button color='secondary' variant='outlined' onClick={handleDelete}>
              Delete
            </Button>
          )}
          <div style={{ flex: '1' }}></div> {/* This will push the next buttons to the right */}
          {addNew ? (
            <Button color='primary' variant='contained' onClick={handleCreate}>
              Save
            </Button>
          ) : (
            <Button color='primary' variant='contained' onClick={handleUpdate}>
              Update
            </Button>
          )}
          <Button color='primary' variant='outlined' onClick={handleClose}>
            Cancel
          </Button>
        </DialogActions>
      </Dialog>
      <ConfirmModal
        title={`${formik.values.name}`}
        message='Are you sure you want to restart this trigger? It may take a few minutes to take effect.'
        onConfirm={() => {
          if (trigger.active) {
            console.log('stop is called with id', formik.values?.ID)
            triggerHandler.handleStop(formik.values?.ID)
          } else {
            triggerHandler.handleRestart(formik.values)
          }
          setOpen(false)
        }}
        show={confirmModal}
        setShow={setConfirmModal}
      />
    </>
  )
}

export default TriggerCard
