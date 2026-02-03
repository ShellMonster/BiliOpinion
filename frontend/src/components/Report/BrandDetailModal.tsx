import React, { useMemo } from 'react';
import Modal from '../common/Modal';
import type { BrandRanking, BrandAnalysis, TypicalComment, Dimension } from '../../types/report';

interface BrandDetailModalProps {
  isOpen: boolean;
  onClose: () => void;
  brandName: string;
  ranking?: BrandRanking;
  analysis?: BrandAnalysis;
  topComments?: TypicalComment[];
  badComments?: TypicalComment[];
  dimensions?: Dimension[];
}

/**
 * 品牌详细信息弹窗组件
 * 
 * 展示品牌的深度分析数据，包括：
 * - 核心指标得分详情
 * - 优劣势详细分析
 * - 典型正面/负面用户评价
 */
export const BrandDetailModal: React.FC<BrandDetailModalProps> = ({
  isOpen,
  onClose,
  brandName,
  ranking,
  analysis,
  topComments = [],
  badComments = [],
  dimensions = []
}) => {
  // 综合得分颜色辅助函数
  const getScoreColor = (score: number) => {
    if (score >= 90) return 'text-emerald-600 bg-emerald-50 border-emerald-100';
    if (score >= 80) return 'text-blue-600 bg-blue-50 border-blue-100';
    if (score >= 70) return 'text-amber-600 bg-amber-50 border-amber-100';
    return 'text-rose-600 bg-rose-50 border-rose-100';
  };

  // 进度条颜色
  const getProgressColor = (score: number) => {
    if (score >= 9) return 'bg-emerald-500';
    if (score >= 8) return 'bg-blue-500';
    if (score >= 7) return 'bg-amber-500';
    return 'bg-rose-500';
  };

  // 匹配维度描述
  const getDimensionDesc = (dimName: string) => {
    const dim = dimensions.find(d => d.name === dimName);
    return dim?.description || '';
  };

  // 排序分数项 (高分在前)
  const sortedScores = useMemo(() => {
    if (!ranking?.scores) return [];
    return Object.entries(ranking.scores)
      .sort(([, a], [, b]) => b - a);
  }, [ranking?.scores]);

  return (
    <Modal
      isOpen={isOpen}
      onClose={onClose}
      title={`${brandName} 详细分析报告`}
    >
      <div className="space-y-8">
        {/* 1. 核心概览区域 */}
        {ranking && (
          <div className="flex flex-col md:flex-row gap-6 items-center bg-slate-50 rounded-2xl p-6 border border-slate-100">
            {/* 左侧：综合得分大卡片 */}
            <div className={`flex flex-col items-center justify-center p-6 rounded-2xl border-2 min-w-[160px] ${getScoreColor(ranking.overall_score)}`}>
              <span className="text-sm font-semibold uppercase tracking-wider opacity-80">综合得分</span>
              <span className="text-5xl font-black mt-2 tracking-tighter">
                {ranking.overall_score.toFixed(1)}
              </span>
              <div className="mt-2 text-xs font-medium px-2 py-1 rounded-full bg-white/50 backdrop-blur-sm">
                排名 #{ranking.rank}
              </div>
            </div>

            {/* 右侧：维度得分网格 */}
            <div className="flex-1 w-full grid grid-cols-1 sm:grid-cols-2 gap-4">
              {sortedScores.map(([name, score]) => (
                <div key={name} className="flex flex-col gap-1.5">
                  <div className="flex justify-between items-end">
                    <span className="text-sm font-medium text-slate-700">{name}</span>
                    <span className="text-sm font-bold text-slate-900">{score.toFixed(1)}</span>
                  </div>
                  <div className="h-2 w-full bg-slate-200 rounded-full overflow-hidden">
                    <div 
                      className={`h-full rounded-full ${getProgressColor(score)} transition-all duration-500`}
                      style={{ width: `${score * 10}%` }}
                    />
                  </div>
                  <p className="text-[10px] text-slate-400 truncate">
                    {getDimensionDesc(name)}
                  </p>
                </div>
              ))}
            </div>
          </div>
        )}

        {/* 2. 优劣势分析 */}
        {analysis && (
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* 优势 */}
            <div className="space-y-3">
              <h3 className="flex items-center text-lg font-bold text-slate-800">
                <span className="w-8 h-8 rounded-lg bg-emerald-100 text-emerald-600 flex items-center justify-center mr-2">
                  <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M5 13l4 4L19 7" />
                  </svg>
                </span>
                核心优势
              </h3>
              <div className="flex flex-wrap gap-2">
                {analysis.strengths && analysis.strengths.length > 0 ? (
                  analysis.strengths.map((s, i) => (
                    <span key={i} className="px-3 py-1.5 rounded-lg bg-emerald-50 text-emerald-700 text-sm font-medium border border-emerald-100">
                      {s}
                    </span>
                  ))
                ) : (
                  <span className="text-slate-400 italic text-sm">暂无明显优势数据</span>
                )}
              </div>
            </div>

            {/* 劣势 */}
            <div className="space-y-3">
              <h3 className="flex items-center text-lg font-bold text-slate-800">
                <span className="w-8 h-8 rounded-lg bg-rose-100 text-rose-600 flex items-center justify-center mr-2">
                  <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2.5} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-1.333-2.694-1.333-3.464 0L3.34 16c-.77 1.333.192 3 1.732 3z" />
                  </svg>
                </span>
                待改进点
              </h3>
              <div className="flex flex-wrap gap-2">
                {analysis.weaknesses && analysis.weaknesses.length > 0 ? (
                  analysis.weaknesses.map((w, i) => (
                    <span key={i} className="px-3 py-1.5 rounded-lg bg-rose-50 text-rose-700 text-sm font-medium border border-rose-100">
                      {w}
                    </span>
                  ))
                ) : (
                  <span className="text-slate-400 italic text-sm">暂无明显劣势数据</span>
                )}
              </div>
            </div>
          </div>
        )}

        {/* 3. 用户评论原声 */}
        <div className="space-y-4">
          <div className="flex items-center justify-between border-b border-slate-100 pb-2">
            <h3 className="text-lg font-bold text-slate-800">用户原声</h3>
            <span className="text-xs font-medium text-slate-400">基于AI提取的典型评论</span>
          </div>
          
          <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
            {/* 好评 */}
            <div className="space-y-3">
              <h4 className="text-sm font-bold text-emerald-600 uppercase tracking-wider flex items-center">
                <div className="w-2 h-2 rounded-full bg-emerald-500 mr-2" />
                典型好评
              </h4>
              <div className="space-y-3">
                {topComments.length > 0 ? (
                  topComments.map((comment, i) => (
                    <div key={i} className="bg-emerald-50/50 p-4 rounded-xl border border-emerald-100 text-sm text-slate-700 leading-relaxed hover:bg-emerald-50 transition-colors">
                      "{comment.content}"
                      {comment.score > 0 && (
                        <div className="mt-2 flex items-center">
                           <span className="text-[10px] font-bold px-1.5 py-0.5 rounded bg-emerald-100 text-emerald-700">
                             评分: {comment.score}
                           </span>
                        </div>
                      )}
                    </div>
                  ))
                ) : (
                  <div className="p-4 rounded-xl bg-slate-50 text-slate-400 text-sm text-center italic">
                    暂无典型好评
                  </div>
                )}
              </div>
            </div>

            {/* 差评 */}
            <div className="space-y-3">
              <h4 className="text-sm font-bold text-rose-600 uppercase tracking-wider flex items-center">
                <div className="w-2 h-2 rounded-full bg-rose-500 mr-2" />
                典型差评
              </h4>
              <div className="space-y-3">
                {badComments.length > 0 ? (
                  badComments.map((comment, i) => (
                    <div key={i} className="bg-rose-50/50 p-4 rounded-xl border border-rose-100 text-sm text-slate-700 leading-relaxed hover:bg-rose-50 transition-colors">
                      "{comment.content}"
                      {comment.score > 0 && (
                        <div className="mt-2 flex items-center">
                           <span className="text-[10px] font-bold px-1.5 py-0.5 rounded bg-rose-100 text-rose-700">
                             评分: {comment.score}
                           </span>
                        </div>
                      )}
                    </div>
                  ))
                ) : (
                  <div className="p-4 rounded-xl bg-slate-50 text-slate-400 text-sm text-center italic">
                    暂无典型差评
                  </div>
                )}
              </div>
            </div>
          </div>
        </div>
      </div>
    </Modal>
  );
};
