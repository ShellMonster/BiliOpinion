import React from 'react'
import { Play, MessageCircle, ExternalLink, Video } from 'lucide-react'
import type { VideoSource } from '../../types/report'

interface VideoSourceListProps {
  videos: VideoSource[]
}

const formatNumber = (num: number): string => {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + '万'
  }
  return num.toString()
}

export const VideoSourceList: React.FC<VideoSourceListProps> = ({ videos }) => {
  if (!videos || videos.length === 0) {
    return null
  }

  return (
    <div className="bg-white rounded-xl shadow-sm border border-gray-200 overflow-hidden">
      <div className="px-6 py-4 border-b border-gray-200">
        <div className="flex items-center gap-2">
          <Video className="w-5 h-5 text-pink-500" />
          <h2 className="text-xl font-bold text-gray-800">数据来源</h2>
        </div>
        <p className="text-sm text-gray-500 mt-1">
          本报告基于以下 {videos.length} 个B站视频的评论分析生成
        </p>
      </div>
      
      <div className="divide-y divide-gray-100">
        {videos.map((video, index) => (
          <a
            key={video.bvid}
            href={`https://www.bilibili.com/video/${video.bvid}`}
            target="_blank"
            rel="noopener noreferrer"
            className="flex items-center gap-4 px-6 py-4 hover:bg-gray-50 transition-colors group"
          >
            <span className="flex-shrink-0 w-8 h-8 rounded-full bg-gray-100 flex items-center justify-center text-sm font-medium text-gray-500">
              {index + 1}
            </span>
            
            <div className="flex-1 min-w-0">
              <h3 className="text-sm font-medium text-gray-900 truncate group-hover:text-pink-600 transition-colors">
                {video.title}
              </h3>
              <p className="text-xs text-gray-500 mt-0.5">
                UP主: {video.author}
              </p>
            </div>
            
            <div className="flex items-center gap-4 flex-shrink-0">
              <div className="flex items-center gap-1 text-xs text-gray-500">
                <Play className="w-3.5 h-3.5" />
                <span>{formatNumber(video.play)}</span>
              </div>
              <div className="flex items-center gap-1 text-xs text-gray-500">
                <MessageCircle className="w-3.5 h-3.5" />
                <span>{formatNumber(video.video_review)}</span>
              </div>
              <ExternalLink className="w-4 h-4 text-gray-400 group-hover:text-pink-500 transition-colors" />
            </div>
          </a>
        ))}
      </div>
    </div>
  )
}

export default VideoSourceList
