interface ReportHeaderProps {
  category: string
  reportId: string | number
  exporting: boolean
  onExport: () => void
  onBack: () => void
}

const ReportHeader = ({
  category,
  reportId,
  exporting,
  onExport,
  onBack
}: ReportHeaderProps) => {
  return (
    <div className="flex justify-between items-end mb-6">
      <div>
        <h1 className="text-3xl font-bold text-gray-800">产品分析报告</h1>
        <p className="text-gray-500">{category} | 报告ID: {reportId}</p>
      </div>
      <div className="flex gap-2">
        <button 
          onClick={onExport}
          disabled={exporting}
          className="px-4 py-2 bg-green-500 text-white rounded-lg hover:bg-green-600 transition-colors font-medium text-sm cursor-pointer disabled:opacity-50 disabled:cursor-not-allowed flex items-center gap-2"
        >
          {exporting ? (
            <>
              <svg className="w-4 h-4 animate-spin" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"></circle>
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
              </svg>
              导出中...
            </>
          ) : (
            <>
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              导出PDF
            </>
          )}
        </button>
        <button 
          onClick={onBack}
          className="px-4 py-2 bg-blue-50 text-blue-600 rounded-lg hover:bg-blue-100 transition-colors font-medium text-sm cursor-pointer"
        >
          查看历史
        </button>
      </div>
    </div>
  )
}

export default ReportHeader
