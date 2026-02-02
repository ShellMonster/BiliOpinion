import { ReportStats } from '../../../types/report'

interface KeyStatsCardsProps {
  stats: ReportStats
  brandCount: number
}

/**
 * 关键统计指标卡片组件
 * 展示：分析视频数、评论总数、对比品牌数
 */
export const KeyStatsCards = ({ stats, brandCount }: KeyStatsCardsProps) => {
  return (
    <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
      {/* 视频总数卡片 */}
      <div className="glass-card p-6 text-center">
        <div className="text-3xl font-bold text-blue-600">{stats.total_videos}</div>
        <div className="text-gray-500 text-sm mt-1">分析视频数</div>
      </div>
      
      {/* 评论总数卡片 */}
      <div className="glass-card p-6 text-center">
        <div className="text-3xl font-bold text-purple-600">{stats.total_comments}</div>
        <div className="text-gray-500 text-sm mt-1">评论总数</div>
      </div>
      
      {/* 品牌数卡片 */}
      <div className="glass-card p-6 text-center">
        <div className="text-3xl font-bold text-green-600">{brandCount}</div>
        <div className="text-gray-500 text-sm mt-1">对比品牌数</div>
      </div>
    </div>
  )
}
