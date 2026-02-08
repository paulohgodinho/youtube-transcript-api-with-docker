import os
from flask import Flask, request, jsonify

from ._api import YouTubeTranscriptApi
from ._errors import (
    YouTubeTranscriptApiException,
    NoTranscriptFound,
    VideoUnavailable,
    TranscriptsDisabled,
    InvalidVideoId,
)


app = Flask(__name__)


@app.route("/")
def health_check():
    """Simple health check endpoint"""
    return jsonify({"status": "ok", "service": "youtube-transcript-api"})


@app.route("/fetch")
def fetch_transcript():
    """Fetch YouTube transcript with default parameters"""
    video_id = request.args.get("video_id")
    
    if not video_id:
        return jsonify({"error": "video_id parameter is required"}), 400
    
    try:
        # Initialize API and fetch transcript with default parameters
        ytt_api = YouTubeTranscriptApi()
        transcript = ytt_api.fetch(video_id, languages=["en"], preserve_formatting=False)
        
        # Convert to raw data format
        transcript_data = transcript.to_raw_data()
        
        # Return response with metadata
        return jsonify({
            "video_id": transcript.video_id,
            "language": transcript.language,
            "language_code": transcript.language_code,
            "is_generated": transcript.is_generated,
            "transcript": transcript_data
        })
    
    except InvalidVideoId:
        return jsonify({"error": "Invalid video ID"}), 400
    except NoTranscriptFound:
        return jsonify({"error": "No transcript found for this video"}), 404
    except VideoUnavailable:
        return jsonify({"error": "Video is unavailable"}), 404
    except TranscriptsDisabled:
        return jsonify({"error": "Transcripts are disabled for this video"}), 404
    except YouTubeTranscriptApiException as e:
        return jsonify({"error": str(e)}), 500
    except Exception as e:
        return jsonify({"error": f"Unexpected error: {str(e)}"}), 500


def run_server():
    """Run the Flask server"""
    port = int(os.environ.get("PORT", 5000))
    host = os.environ.get("HOST", "0.0.0.0")
    debug = os.environ.get("DEBUG", "False").lower() == "true"
    
    app.run(host=host, port=port, debug=debug)


if __name__ == "__main__":
    run_server()
