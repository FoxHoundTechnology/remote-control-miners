import Divider from '@mui/material/Divider'
import { styled } from '@mui/material/styles'
import Typography, { TypographyProps } from '@mui/material/Typography'
import MuiListSubheader, { ListSubheaderProps } from '@mui/material/ListSubheader'

import { NavSectionTitle } from 'src/@core/layouts/types'
import { Settings } from 'src/@core/context/settingsContext'

import Translations from 'src/layouts/components/Translations'

interface Props {
  navHover: boolean
  settings: Settings
  item: NavSectionTitle
  collapsedNavWidth: number
  navigationBorderWidth: number
}

const ListSubheader = styled((props: ListSubheaderProps) => <MuiListSubheader component='li' {...props} />)(
  ({ theme }) => ({
    lineHeight: 1,
    display: 'flex',
    position: 'static',
    marginTop: theme.spacing(7),
    marginBottom: theme.spacing(2),
    backgroundColor: 'transparent'
  })
)

const TypographyHeaderText = styled(Typography)<TypographyProps>(({ theme }) => ({
  fontSize: '0.75rem',
  lineHeight: 'normal',
  letterSpacing: '0.21px',
  textTransform: 'uppercase',
  fontWeight: theme.typography.fontWeightMedium
}))

const VerticalNavSectionTitle = (props: Props) => {
  const { item, navHover, settings, collapsedNavWidth, navigationBorderWidth } = props

  const { navCollapsed } = settings

  return (
    <ListSubheader
      className='nav-section-title'
      sx={{
        ...(navCollapsed && !navHover
          ? {
              py: 3.5,
              pr: (collapsedNavWidth - navigationBorderWidth - 24) / 8 - 1,
              pl: (collapsedNavWidth - navigationBorderWidth - 24) / 8 + 0.25
            }
          : { px: 0, py: 1.75 })
      }}
    >
      <Divider
        textAlign='left'
        sx={{
          m: '0 !important',
          lineHeight: 'normal',
          ...(navCollapsed && !navHover
            ? {
                width: 22,
                borderColor: theme => `rgba(${theme.palette.customColors.main}, 0.3)`
              }
            : {
                width: '100%',
                textTransform: 'uppercase',
                '&:before, &:after': { top: 7, transform: 'none' },
                '& .MuiDivider-wrapper': { px: 2.5, fontSize: '0.75rem', letterSpacing: '0.21px' }
              })
        }}
      >
        {navCollapsed && !navHover ? null : (
          <TypographyHeaderText noWrap sx={{ color: 'text.disabled' }}>
            <Translations text={item.sectionTitle} />
          </TypographyHeaderText>
        )}
      </Divider>
    </ListSubheader>
  )
}

export default VerticalNavSectionTitle
