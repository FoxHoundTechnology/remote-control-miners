import { Ref, useState, forwardRef, ReactElement, SetStateAction, Dispatch } from 'react'

import Dialog from '@mui/material/Dialog'
import Button from '@mui/material/Button'
import MenuItem from '@mui/material/MenuItem'
import TextField from '@mui/material/TextField'
import IconButton from '@mui/material/IconButton'
import Typography from '@mui/material/Typography'
import InputLabel from '@mui/material/InputLabel'
import FormControl from '@mui/material/FormControl'
import CardContent from '@mui/material/CardContent'
import Fade, { FadeProps } from '@mui/material/Fade'
import DialogContent from '@mui/material/DialogContent'
import DialogActions from '@mui/material/DialogActions'
import FormControlLabel from '@mui/material/FormControlLabel'
import Select, { SelectChangeEvent } from '@mui/material/Select'

import Icon from 'src/@core/components/icon'

type ConfirmModalProps = {
  title: string
  message?: string
  onConfirm: () => void
  show: boolean
  setShow: (show: boolean) => void
}

const ConfirmModal = ({ title, message, onConfirm, show, setShow }: ConfirmModalProps) => {
  return (
    <Dialog
      fullWidth
      open={show}
      maxWidth='sm'
      scroll='body'
      onClose={() => {
        setShow(false)
      }}
      // TransitionComponent={Transition}
    >
      <DialogContent
        sx={{
          position: 'relative',
          pb: theme => `${theme.spacing(8)} !important`,
          px: theme => [`${theme.spacing(5)} !important`, `${theme.spacing(15)} !important`],
          pt: theme => [`${theme.spacing(8)} !important`, `${theme.spacing(12.5)} !important`]
        }}
      >
        <Typography variant='h5' sx={{ mb: 4 }}>
          Remote Control : {title}
        </Typography>
        <Typography sx={{ mb: 3 }}>{message}</Typography>
      </DialogContent>
      <DialogActions
        sx={{
          justifyContent: 'center',
          px: theme => [`${theme.spacing(5)} !important`, `${theme.spacing(15)} !important`],
          pb: theme => [`${theme.spacing(8)} !important`, `${theme.spacing(12.5)} !important`]
        }}
      >
        <Button
          variant='contained'
          sx={{ mr: 1 }}
          onClick={() => {
            onConfirm()
            setShow(false)
          }}
        >
          Confirm
        </Button>
        <Button
          variant='outlined'
          color='secondary'
          onClick={() => {
            setShow(false)
          }}
        >
          Cancel
        </Button>
      </DialogActions>
    </Dialog>
  )
}

export default ConfirmModal
