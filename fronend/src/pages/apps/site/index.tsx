import React from 'react'
import { DndProvider } from 'react-dnd'
import { HTML5Backend } from 'react-dnd-html5-backend'
import DataCenterSitemapContent from 'src/views/apps/site/DataCenterSitemapContent'

const DataCenterSitemapPage: React.FC = () => {
  return (
    <DndProvider backend={HTML5Backend}>
      <DataCenterSitemapContent />
    </DndProvider>
  )
}

export default DataCenterSitemapPage
