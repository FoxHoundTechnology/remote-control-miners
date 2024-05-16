import Link from 'next/link'

import Box from '@mui/material/Box'
import Grid from '@mui/material/Grid'
import Card from '@mui/material/Card'
import Switch from '@mui/material/Switch'
import Button from '@mui/material/Button'
import Typography from '@mui/material/Typography'
import CardContent from '@mui/material/CardContent'

import Icon from 'src/@core/components/icon'
import ShowChartIcon from '@mui/icons-material/ShowChart'

import { ComponentType, useState } from 'react'
import { Divider, FormControl, FormControlLabel, SvgIconProps, TextField, useTheme } from '@mui/material'

interface NotificationType {
  title: string
  icon: ComponentType<SvgIconProps>
  enabled: boolean
  subtitle: string
}

/*
const connectedAccountsArr: NotificationType[] = [
  {
    enabled: true,
    title: 'Google',
    icon: ,
    subtitle: 'Calendar and Contacts'
  },
  {
    enabled: false,
    title: 'Slack',
    icon ,
    subtitle: 'Communications'
  },
  {
    enabled: true,
    title: 'Github',
    icon:,
    subtitle: 'Manage your Git repositories'
  },
  {
    enabled: true,
    title: 'Mailchimp',
    icon:, 
    subtitle: 'Email marketing service',
  },
  {
    title: 'Asana',
    enabled: false,
    icon:,
    subtitle: 'Communication',
  }
]
*/

const NotificationSettings = () => {
  const theme = useTheme()

  const initialState = {
    alertHashRateEnabled: false,
    alertBelowHashRate: 0,
    warnHashRateEnabled: false,
    warnBelowHashRate: 0,
    alertFanEnabled: false,
    alertAboveFan: 0,
    warnFanEnabled: false,
    warnAboveFan: 0,
    alertTempEnabled: false,
    alertAboveTemp: 0,
    warnTempEnabled: false,
    warnAboveTemp: 0
  }

  const [notificationSettings, setNotificationSettings] = useState(initialState)

  const handleInputChange = (keyVal: string, value: string | React.ChangeEvent<HTMLInputElement>) => {
    if (typeof value === 'string' || typeof value === 'number') {
      setNotificationSettings(prev => ({ ...prev, [keyVal]: value }))
    } else {
      const inputValue = (value.target as HTMLInputElement).value
      setNotificationSettings(prev => ({ ...prev, [keyVal]: inputValue }))
    }
  }

  return (
    <Grid container spacing={6}>
      {/* Connected Accounts Cards */}
      <Grid item xs={12}>
        <Card>
          <CardContent>
            <Box sx={{ mb: 5 }}>
              <Typography sx={{ fontWeight: 500 }}>Alerts and Warning</Typography>
              <Typography variant='body2'>Display content from your connected accounts on your site</Typography>
            </Box>

            <Grid container spacing={5}>
              <Grid item xs={6} md={3}>
                <FormControl
                  required
                  fullWidth
                  variant='outlined'
                  onChange={e => handleInputChange('alertHashrateEnabled', e as React.ChangeEvent<HTMLInputElement>)}
                >
                  <FormControlLabel
                    control={<Switch color='primary' value={notificationSettings?.alertHashRateEnabled} />}
                    label='Hashrate Alert'
                    labelPlacement='end'
                  />
                </FormControl>
              </Grid>

              <Grid item xs={6} md={3}>
                <FormControl
                  required
                  fullWidth
                  variant='outlined'
                  onChange={e => handleInputChange('alertBelowHashRate', e as React.ChangeEvent<HTMLInputElement>)}
                >
                  <TextField
                    type='number'
                    placeholder='TH/s (Alert)'
                    value={notificationSettings?.alertBelowHashRate}
                  />
                </FormControl>
              </Grid>

              <Grid item xs={6} md={3}>
                <FormControl
                  required
                  fullWidth
                  variant='outlined'
                  onChange={e => handleInputChange('warnHashrateEnabled', e as React.ChangeEvent<HTMLInputElement>)}
                >
                  <FormControlLabel
                    control={<Switch color='primary' value={notificationSettings?.warnHashRateEnabled} />}
                    label='Hashrate Warning'
                    labelPlacement='end'
                  />
                </FormControl>
              </Grid>

              <Grid item xs={6} md={3}>
                <FormControl
                  required
                  fullWidth
                  variant='outlined'
                  onChange={e => handleInputChange('warnBelowHashRate', e as React.ChangeEvent<HTMLInputElement>)}
                >
                  <TextField
                    type='number'
                    placeholder='TH/s (Warning)'
                    value={notificationSettings?.warnBelowHashRate}
                  />
                </FormControl>
              </Grid>
              <Grid item xs={12}>
                <Divider />
              </Grid>

              <Grid item xs={6} md={3}>
                <FormControl
                  required
                  fullWidth
                  variant='outlined'
                  onChange={e => handleInputChange('alertFanEneabled', e as React.ChangeEvent<HTMLInputElement>)}
                >
                  <FormControlLabel
                    control={<Switch color='primary' value={notificationSettings?.alertFanEnabled} />}
                    label='Fan Speed Alert'
                    labelPlacement='end'
                  />
                </FormControl>
              </Grid>

              <Grid item xs={6} md={3}>
                <FormControl
                  required
                  fullWidth
                  variant='outlined'
                  onChange={e => handleInputChange('alertAboveFan', e as React.ChangeEvent<HTMLInputElement>)}
                >
                  <TextField type='number' placeholder='BPM (Alert)' value={notificationSettings?.alertAboveFan} />
                </FormControl>
              </Grid>

              <Grid item xs={6} md={3}>
                <FormControl
                  required
                  fullWidth
                  variant='outlined'
                  onChange={e => handleInputChange('warnFanEnabled', e as React.ChangeEvent<HTMLInputElement>)}
                >
                  <FormControlLabel
                    control={<Switch color='primary' value={notificationSettings?.warnFanEnabled} />}
                    label='Fan Speed Warning'
                    labelPlacement='end'
                  />
                </FormControl>
              </Grid>

              <Grid item xs={6} md={3}>
                <FormControl
                  required
                  fullWidth
                  variant='outlined'
                  onChange={e => handleInputChange('warnAboveFan', e as React.ChangeEvent<HTMLInputElement>)}
                >
                  <TextField type='number' placeholder='BPM (Warning)' value={notificationSettings?.warnAboveFan} />
                </FormControl>
              </Grid>
              <Grid item xs={12}>
                <Divider />
              </Grid>

              <Grid item xs={6} md={3}>
                <FormControl
                  required
                  fullWidth
                  variant='outlined'
                  onChange={e => handleInputChange('alertTempEnabled', e as React.ChangeEvent<HTMLInputElement>)}
                >
                  <FormControlLabel
                    control={<Switch color='primary' value={notificationSettings?.alertTempEnabled} />}
                    label='Temperature Alert'
                    labelPlacement='end'
                  />
                </FormControl>
              </Grid>

              <Grid item xs={6} md={3}>
                <FormControl
                  required
                  fullWidth
                  variant='outlined'
                  onChange={e => handleInputChange('alertAboveTemp', e as React.ChangeEvent<HTMLInputElement>)}
                >
                  <TextField type='number' placeholder='°C (Alert)' value={notificationSettings?.alertAboveTemp} />
                </FormControl>
              </Grid>

              <Grid item xs={6} md={3}>
                <FormControl
                  required
                  fullWidth
                  variant='outlined'
                  onChange={e => handleInputChange('warnTempEnabled', e as React.ChangeEvent<HTMLInputElement>)}
                >
                  <FormControlLabel
                    control={<Switch color='primary' value={notificationSettings?.warnTempEnabled} />}
                    label='Temperature Warning'
                    labelPlacement='end'
                  />
                </FormControl>
              </Grid>

              <Grid item xs={6} md={3}>
                <FormControl
                  required
                  fullWidth
                  variant='outlined'
                  onChange={e => handleInputChange('warnAboveTemp', e as React.ChangeEvent<HTMLInputElement>)}
                >
                  <TextField type='number' placeholder='°C (Warning)' value={notificationSettings?.warnAboveTemp} />
                </FormControl>
              </Grid>
              <Grid item xs={12}>
                <Box display='flex' justifyContent='flex-end'>
                  <Button variant='contained' sx={{ mr: 4 }} onClick={() => {}}>
                    Save
                  </Button>
                  <Button color='info' variant='outlined' onClick={() => {}}>
                    Reset
                  </Button>
                </Box>
              </Grid>
            </Grid>

            {/*connectedAccountsArr.map(account => {
              return (
                <Box
                  key={account.title}
                  sx={{
                    gap: 2,
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'space-between',
                    '&:not(:last-of-type)': { mb: 4 }
                  }}
                >
                  <Box sx={{ display: 'flex', alignItems: 'center' }}>
                    <Box sx={{ mr: 4, display: 'flex', justifyContent: 'center' }}>
                      <img src={account.logo} alt={account.title} height='36' width='36' />
                    </Box>
                    <div>
                      <Typography sx={{ fontWeight: 600 }}>{account.title}</Typography>
                      <Typography variant='body2'>{account.subtitle}</Typography>
                    </div>
                  </Box>
                  <Switch defaultChecked={account.checked} />
                </Box>
              )
            })*/}
          </CardContent>
        </Card>
      </Grid>
    </Grid>
  )
}

export default NotificationSettings
