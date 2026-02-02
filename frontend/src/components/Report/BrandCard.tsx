import React from 'react';
import type { BrandRanking, BrandAnalysis } from '../../types/report';

interface BrandCardProps {
  ranking: BrandRanking;
  analysis?: BrandAnalysis;
  onClick: () => void;
}

/**
 * 品牌展示卡片组件
 * 
 * 展示品牌的核心信息，包括：
 * - 排名 (Rank)
 * - 品牌名称 (Brand Name)
 * - 综合得分 (Overall Score)
 * - 优劣势标签 (Strengths/Weaknesses)
 */
export const BrandCard: React.FC<BrandCardProps> = ({ ranking, analysis, onClick }) => {
  // 根据分数获取颜色类名
  const getScoreColor = (score: number) => {
    if (score >= 90) return 'text-emerald-500';
    if (score >= 80) return 'text-blue-500';
    if (score >= 70) return 'text-amber-500';
    return 'text-rose-500';
  };

  // 根据排名获取背景样式
  const getRankStyle = (rank: number) => {
    switch (rank) {
      case 1:
        return 'bg-gradient-to-br from-yellow-300 to-amber-500 text-white border-none shadow-amber-200/50';
      case 2:
        return 'bg-gradient-to-br from-slate-300 to-slate-400 text-white border-none shadow-slate-200/50';
      case 3:
        return 'bg-gradient-to-br from-orange-300 to-orange-400 text-white border-none shadow-orange-200/50';
      default:
        return 'bg-slate-100 text-slate-600 border border-slate-200';
    }
  };

  return (
    <div 
      className="group relative overflow-hidden rounded-2xl bg-white border border-slate-100 p-6 shadow-sm transition-all duration-300 hover:shadow-xl hover:-translate-y-1 cursor-pointer"
      onClick={onClick}
    >
      {/* 装饰性背景光晕 */}
      <div className="absolute top-0 right-0 -mr-16 -mt-16 h-48 w-48 rounded-full bg-gradient-to-br from-blue-50 to-purple-50 opacity-50 blur-3xl transition-all duration-500 group-hover:scale-150" />

      <div className="relative z-10 flex flex-col gap-4">
        {/* 头部：排名与分数 */}
        <div className="flex items-center justify-between">
          <div className={`flex h-10 w-10 items-center justify-center rounded-xl text-lg font-bold shadow-sm ${getRankStyle(ranking.rank)}`}>
            #{ranking.rank}
          </div>
          <div className="flex flex-col items-end">
            <span className={`text-3xl font-black tracking-tight ${getScoreColor(ranking.overall_score)}`}>
              {ranking.overall_score.toFixed(1)}
            </span>
            <span className="text-xs font-medium text-slate-400 uppercase tracking-wider">综合得分</span>
          </div>
        </div>

        {/* 品牌名称 */}
        <div>
          <h3 className="text-2xl font-bold text-slate-800 tracking-tight group-hover:text-blue-600 transition-colors">
            {ranking.brand}
          </h3>
        </div>

        {/* 优劣势标签 */}
        {analysis && (
          <div className="mt-2 space-y-3">
            {/* 优势 */}
            {analysis.strengths.length > 0 && (
              <div className="flex flex-wrap gap-1.5">
                {analysis.strengths.slice(0, 3).map((strength, index) => (
                  <span 
                    key={`strength-${index}`} 
                    className="inline-flex items-center px-2 py-0.5 rounded-md bg-emerald-50 text-emerald-700 text-xs font-medium border border-emerald-100"
                  >
                    <svg className="w-3 h-3 mr-1 text-emerald-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                    {strength}
                  </span>
                ))}
              </div>
            )}
            
            {/* 劣势 (仅显示一个关键劣势，避免负面情绪过重) */}
            {analysis.weaknesses.length > 0 && (
              <div className="flex flex-wrap gap-1.5">
                {analysis.weaknesses.slice(0, 1).map((weakness, index) => (
                  <span 
                    key={`weakness-${index}`} 
                    className="inline-flex items-center px-2 py-0.5 rounded-md bg-rose-50 text-rose-700 text-xs font-medium border border-rose-100"
                  >
                    <svg className="w-3 h-3 mr-1 text-rose-500" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                    </svg>
                    {weakness}
                  </span>
                ))}
              </div>
            )}
          </div>
        )}
        
        {/* 查看详情提示 */}
        <div className="mt-2 flex items-center text-sm font-medium text-blue-600 opacity-0 transform translate-y-2 transition-all duration-300 group-hover:opacity-100 group-hover:translate-y-0">
          查看详细分析 
          <svg className="w-4 h-4 ml-1" fill="none" viewBox="0 0 24 24" stroke="currentColor">
            <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17 8l4 4m0 0l-4 4m4-4H3" />
          </svg>
        </div>
      </div>
    </div>
  );
};
