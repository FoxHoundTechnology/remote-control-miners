import { FormHelperText, IconProps } from '@mui/material'
import Box from '@mui/material/Box'
import Grid, { GridProps } from '@mui/material/Grid'
import Radio from '@mui/material/Radio'
import Typography from '@mui/material/Typography'
import { ChangeEvent } from 'react'

import { CustomRadioIconsData, CustomRadioIconsProps } from 'src/@core/components/custom-radio/types'

import Icon from 'src/@core/components/icon'
import { ThemeColor } from 'src/@core/layouts/types'

export type RadioIconsProps = {
  name: string
  icon?: string
  selected: string
  color?: ThemeColor
  gridProps: GridProps
  data: CustomRadioIconsData
  iconProps?: Omit<IconProps, 'icon'>
  handleChange: any
  helperText?: string
  enabled?: boolean
}

const CustomRadioIcons = (props: RadioIconsProps) => {
  const {
    data,
    icon,
    name,
    selected,
    gridProps,
    iconProps,
    handleChange,
    color = 'primary',
    helperText,
    enabled
  } = props

  const { title, value, content } = data

  const renderComponent = () => {
    return (
      <Grid item {...gridProps}>
        <Box
          onClick={() => {
            if (enabled) handleChange(value)
          }}
          sx={{
            p: 4,
            height: '100%',
            display: 'flex',
            borderRadius: 1,
            cursor: 'pointer',
            position: 'relative',
            alignItems: 'center',
            flexDirection: 'column',
            border: theme => `1px solid ${theme.palette.divider}`,
            ...(selected === value
              ? { borderColor: `${color}.main` }
              : { '&:hover': { borderColor: theme => `rgba(${theme.palette.customColors.main}, 0.25)` } })
          }}
        >
          {/*icon ? <Icon icon={icon} {...iconProps} /> : null */}
          {title ? (
            typeof title === 'string' ? (
              <Typography
                sx={{
                  fontWeight: 700,
                  ...(content ? { mb: 1 } : { my: 'auto' })
                }}
              >
                {title}
              </Typography>
            ) : (
              title
            )
          ) : null}
          <Radio
            name={name}
            size='small'
            color={color}
            value={value}
            onChange={handleChange}
            checked={selected === value}
            sx={{ mb: -2, ...(!icon && !title && !content && { mt: -2 }) }}
          />
        </Box>
      </Grid>
    )
  }

  return data ? renderComponent() : null
}

export default CustomRadioIcons
