'use client';

import { useState } from 'react';
import { io } from 'socket.io-client';
import { ToastContainer, toast } from 'react-toastify';
import 'react-toastify/dist/ReactToastify.css';
import logoImage from '../public/logo.png'
import Image from 'next/image';
import { FaVideo } from 'react-icons/fa6';
import { FaMusic } from 'react-icons/fa';

export default function UploadPage() {
  const [file, setFile] = useState<File | null>(null);
  const [outputFormat, setOutputFormat] = useState('mp4');
  const [quality, setQuality] = useState('low');
  const [progress, setProgress] = useState<number>(0);
  const [operationId, setOperationId] = useState<string | null>(null);
  const [downloadUrl, setDownloadUrl] = useState<string | null>(null);
  const [activeTab, setActiveTab] = useState('video'); // Video ou audio

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    if (event.target.files && event.target.files.length > 0) {
      setFile(event.target.files[0]);
    }
  };

  const showError = (message: string) => {
    toast.error(message, {
      position: 'top-right',
      autoClose: 5000,
      hideProgressBar: false,
      closeOnClick: true,
      pauseOnHover: true,
      draggable: true,
      theme: 'dark',
    });
  };

  const showSuccess = (message: string) => {
    toast.success(message, {
      position: 'top-right',
      autoClose: 3000,
      hideProgressBar: false,
      closeOnClick: true,
      pauseOnHover: true,
      draggable: true,
      theme: 'dark',
    });
  };

  const handleUpload = async () => {
    if (!file) {
      showError('Selecione um arquivo primeiro!');
      return;
    }
    setProgress(0);
    const formData = new FormData();
    formData.append('file', file);
    formData.append('outputFormat', outputFormat);
    formData.append('quality', quality);

    try {
      const response = await fetch('http://localhost:3000/conversion/upload', {
        method: 'POST',
        body: formData,
      });

      const data = await response.json();

      if (!response.ok) {
        throw new Error(data.message || 'Erro ao enviar o arquivo.');
      }

      setOperationId(data.operationId);
      listenToWebSocket(data.operationId);
    } catch (error: any) {
      console.error('Erro ao enviar o arquivo:', error);
      showError(error.message || 'Erro ao enviar o arquivo.');
    }
  };

  const listenToWebSocket = (operationId: string) => {
    const socket = io('http://localhost:3000', {
      path: '/ws/',
      query: { operationId },
    });

    socket.on('jobUpdate', (update) => {
      if (update.status === 'progress') {
        setProgress(update.progress || 0);
      } else if (update.status === 'success') {
        setProgress(100);
        setDownloadUrl(`http://localhost:3000/conversion/download/${update.new_file_name}`);
        showSuccess('Conversão concluída com sucesso!');
        socket.disconnect();
      } else if (update.status === 'error') {
        showError('Erro na conversão');
        socket.disconnect();
      }
    });
  };

  return (
    <div className="min-h-screen flex flex-col items-center justify-center bg-black text-white p-4">
      <ToastContainer />
      <Image src={logoImage} alt='logo' width={240} />

      <div className="bg-gray-900 p-6 rounded-lg shadow-lg w-96 text-center">
        {!downloadUrl && (
          <>
            <input
              type="file"
              onChange={handleFileChange}
              className="mb-4 w-full p-2 bg-gray-700 rounded-lg text-white"
            />

            <div className="tabs w-full">
              <div className="tab-list flex">
                <button
                  onClick={() => {
                    setActiveTab('video');
                    setOutputFormat('mp4');
                    setQuality('low');
                  }}
                  className={`tab-button py-2 w-full rounded-t-md transition-all px-4 flex items-center justify-center
      ${activeTab === 'video' ? 'bg-gray-400' : 'bg-gray-700'} 
      ${activeTab === 'video' ? 'text-black' : 'text-white'}`}
                >
                  <FaVideo className="mr-2" /> {/* Adiciona margem à direita do ícone */}
                  Video
                </button>
                <button
                  onClick={() => {
                    setActiveTab('audio');
                    setOutputFormat('mp3');
                    setQuality('low');
                  }}
                  className={`tab-button py-2 w-full rounded-t-md transition-all px-4 flex items-center justify-center
      ${activeTab === 'audio' ? 'bg-gray-400' : 'bg-gray-700'} 
      ${activeTab === 'audio' ? 'text-black' : 'text-white'}`}
                >
                  <FaMusic className="mr-2" /> {/* Adiciona margem à direita do ícone */}
                  Áudio
                </button>
              </div>
            </div>

            <div className={`mb-4 p-4 ${activeTab ? 'bg-gray-400' : 'bg-gray-700'} ${activeTab ? 'text-black' : 'bg-gray-700'}`}>
              <label className="block text-left mb-2">Formato</label>
              <select
                onChange={(e) => setOutputFormat(e.target.value)}
                value={outputFormat}
                className="p-2 rounded-lg bg-gray-700 text-white w-full"
              >
                {activeTab === 'video' && (
                  <>
                    <option value="mp4">MP4</option>
                    <option value="avi">AVI</option>
                    <option value="mkv">MKV</option>
                    <option value="mov">MOV</option>
                    <option value="flv">FLV</option>
                    <option value="webm">WEBM</option>
                    <option value="gif">GIF</option>
                  </>
                )}
                {activeTab === 'audio' && (
                  <>
                    <option value="mp3">MP3</option>
                    <option value="wav">WAV</option>
                    <option value="aac">AAC</option>
                    <option value="flac">FLAC</option>
                    <option value="wma">WMA</option>
                  </>
                )}
              </select>

              <label className="block text-left mb-2">Qualidade</label>
              <select
                onChange={(e) => setQuality(e.target.value)}
                value={quality}
                className="p-2 rounded-lg bg-gray-700 text-white w-full"
              >
                <option value="low">Baixa</option>
                <option value="medium">Média</option>
                <option value="high">Alta</option>
              </select>
            </div>

            <button
              onClick={handleUpload}
              className="w-full p-3 bg-blue-600 hover:bg-blue-500 text-white font-bold rounded-lg transition"
            >
              Enviar
            </button>
          </>
        )}

        {operationId && (
          <div className="mt-4">
            <h3 className="text-lg font-semibold">Progresso:</h3>
            <div className="w-full bg-gray-700 rounded-full h-4 mt-2">
              <div
                className="bg-blue-500 h-4 rounded-full"
                style={{ width: `${progress}%` }}
              ></div>
            </div>
            <p className="text-sm mt-2">{progress}%</p>
          </div>
        )}

        {downloadUrl && (
          <div className="flex flex-col items-center mt-4">
            <a
              href={downloadUrl}
              download
              className="mb-4 p-3 bg-green-600 hover:bg-green-500 text-white font-bold rounded-lg transition"
            >
              Baixar Arquivo
            </a>
            <button
              onClick={() => {
                setFile(null);
                setOperationId(null);
                setDownloadUrl(null);
                setProgress(0);
              }}
              className="p-3 bg-gray-600 hover:bg-gray-500 text-white font-bold rounded-lg transition"
            >
              Fazer outra conversão
            </button>
          </div>
        )}
      </div>
    </div >
  );
}

