import React, { useMemo } from 'react';
import ReactECharts from 'echarts-for-react';
import type { ReportData } from '../../../types/report';

interface BrandScoreChartProps {
  data: ReportData;
}

/**
 * 品牌综合得分排名柱状图组件
 * 只显示前10名品牌，避免过于拥挤
 * @param data 报告数据
 */
export const BrandScoreChart: React.FC<BrandScoreChartProps> = ({ data }) => {
  // 品牌颜色配置
  const colors = ['#3b82f6', '#8b5cf6', '#ec4899', '#f59e0b', '#10b981'];

  // 只取前10名品牌，避免图表过于拥挤
  const topRankings = useMemo(() =>
    [...data.rankings].slice(0, 10).reverse(),
    [data.rankings]
  );

  // 使用 useMemo 优化图表配置计算
  const option = useMemo(() => ({
    tooltip: {
      trigger: 'axis',
      axisPointer: {
        type: 'shadow'
      }
    },
    grid: {
      left: 80,
      right: 40,
      top: 10,
      bottom: 10
    },
    xAxis: {
      type: 'value',
      max: 10,
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        color: '#6b7280'
      },
      splitLine: {
        lineStyle: {
          color: '#f0f0f0'
        }
      }
    },
    yAxis: {
      type: 'category',
      data: topRankings.map(r => r.brand),
      axisLine: {
        show: false
      },
      axisTick: {
        show: false
      },
      axisLabel: {
        color: '#6b7280',
        fontSize: 12
      }
    },
    series: [{
      type: 'bar',
      data: topRankings.map((r, index) => ({
        value: Math.round(r.overall_score * 10) / 10,
        itemStyle: {
          color: colors[index % colors.length],
          borderRadius: [0, 8, 8, 0]
        }
      })),
      barWidth: 20,
      label: {
        show: true,
        position: 'right',
        formatter: '{c}',
        color: '#374151',
        fontWeight: 'bold',
        fontSize: 12
      }
    }]
  }), [topRankings, colors]);

  return (
    <div className="w-full bg-white rounded-lg p-4 shadow-sm border border-gray-200">
      <h3 className="text-lg font-semibold text-gray-800 mb-4">品牌综合得分 (Top 10)</h3>
      <ReactECharts option={option} style={{ height: '280px' }} />
    </div>
  );
};
