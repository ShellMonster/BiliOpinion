import React, { useMemo } from 'react';
import ReactECharts from 'echarts-for-react';
import type { ReportData } from '../../../types/report';

interface BrandRadarChartProps {
  data: ReportData;
}

/**
 * 品牌多维度对比雷达图组件
 * @param data 报告数据
 */
export const BrandRadarChart: React.FC<BrandRadarChartProps> = ({ data }) => {
  // 品牌颜色配置
  const colors = ['#3b82f6', '#8b5cf6', '#ec4899', '#f59e0b', '#10b981'];

  // 使用 useMemo 优化图表配置计算
  const option = useMemo(() => ({
    tooltip: {
      trigger: 'item'
    },
    legend: {
      data: data.brands.slice(0, 3),
      bottom: 10,
      textStyle: {
        fontSize: 12
      }
    },
    radar: {
      indicator: data.dimensions.map(dim => ({
        name: dim.name,
        max: 100
      })),
      splitNumber: 4,
      radius: '65%', // 减小雷达图半径，给外部标签留更多空间
      center: ['50%', '50%'],
      name: {
        textStyle: {
          color: '#6b7280',
          fontSize: 10 // 减小字体
        },
        formatter: (value: string) => {
          // 长文本换行处理，每行最多4个字符
          if (value.length > 4) {
            return value.slice(0, 4) + '\n' + value.slice(4);
          }
          return value;
        }
      },
      splitLine: {
        lineStyle: {
          color: '#e5e7eb'
        }
      },
      splitArea: {
        show: true,
        areaStyle: {
          color: ['rgba(255, 255, 255, 0)', 'rgba(249, 250, 251, 0.5)']
        }
      },
      axisLine: {
        lineStyle: {
          color: '#e5e7eb'
        }
      }
    },
    series: [{
      type: 'radar',
      data: data.brands.slice(0, 3).map((brand, index) => ({
        value: data.dimensions.map(dim => 
          data.scores[brand]?.[dim.name] ? data.scores[brand][dim.name] * 10 : 0
        ),
        name: brand,
        lineStyle: {
          color: colors[index],
          width: 2
        },
        areaStyle: {
          color: colors[index],
          opacity: 0.2
        },
        itemStyle: {
          color: colors[index]
        }
      }))
    }]
  }), [data, colors]);

  return (
    <div className="w-full bg-white rounded-lg p-4 shadow-sm border border-gray-200">
      <h3 className="text-lg font-semibold text-gray-800 mb-4">品牌维度对比</h3>
      <ReactECharts option={option} style={{ height: '320px' }} />
    </div>
  );
};
