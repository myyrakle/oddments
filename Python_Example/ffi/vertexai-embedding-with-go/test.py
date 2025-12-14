"""
Test script for Vertex AI Embedding Library

This script demonstrates how to use the Go-based Vertex AI embedding library from Python.
"""

import os
import sys
from vertexai_client import VertexAIEmbedding


def test_text_embedding():
    """Test text embedding functionality."""
    print("=" * 60)
    print("Testing Text Embedding")
    print("=" * 60)

    # Get project ID from environment or use default
    project_id = os.environ.get("GCP_PROJECT_ID", "your-project-id")
    location = os.environ.get("GCP_LOCATION", "us-central1")
    credentials_path = os.environ.get("GOOGLE_APPLICATION_CREDENTIALS")

    print(f"Project ID: {project_id}")
    print(f"Location: {location}")
    print(f"Credentials: {credentials_path or 'Using default credentials'}")
    print()

    # Initialize client
    try:
        client = VertexAIEmbedding(
            project_id=project_id, location=location, credentials_json=credentials_path
        )
        print("✓ Client initialized successfully")
    except Exception as e:
        print(f"✗ Failed to initialize client: {e}")
        return False

    # Test texts
    test_texts = [
        "Hello, world!",
        "The quick brown fox jumps over the lazy dog.",
        "Machine learning is fascinating.",
        "Python and Go are great programming languages.",
    ]

    print("\nTesting text embeddings:")
    print("-" * 60)

    for i, text in enumerate(test_texts, 1):
        try:
            embedding = client.get_text_embedding(text)
            print(f"\n{i}. Text: {text}")
            print(f"   Embedding dimension: {len(embedding)}")
            print(f"   First 5 values: {[f'{x:.4f}' for x in embedding[:5]]}")
            print(f"   Last 5 values: {[f'{x:.4f}' for x in embedding[-5:]]}")
            print("   ✓ Success")
        except Exception as e:
            print(f"\n{i}. Text: {text}")
            print(f"   ✗ Error: {e}")
            return False

    return True


def test_image_embedding():
    """Test image embedding functionality."""
    print("\n" + "=" * 60)
    print("Testing Image Embedding")
    print("=" * 60)

    # Get project ID from environment
    project_id = os.environ.get("GCP_PROJECT_ID", "your-project-id")
    location = os.environ.get("GCP_LOCATION", "us-central1")
    credentials_path = os.environ.get("GOOGLE_APPLICATION_CREDENTIALS")

    # Check if test image exists
    test_image_path = os.environ.get("TEST_IMAGE_PATH")

    if not test_image_path or not os.path.exists(test_image_path):
        print("⚠ Skipping image embedding test")
        print("  Set TEST_IMAGE_PATH environment variable to test image embedding")
        return True

    print(f"Test image: {test_image_path}")

    # Initialize client
    try:
        client = VertexAIEmbedding(
            project_id=project_id, location=location, credentials_json=credentials_path
        )
        print("✓ Client initialized successfully")
    except Exception as e:
        print(f"✗ Failed to initialize client: {e}")
        return False

    # Test image embedding
    print("\nTesting image embedding:")
    print("-" * 60)

    try:
        embedding = client.get_image_embedding(image_path=test_image_path)
        print(f"\nEmbedding dimension: {len(embedding)}")
        print(f"First 5 values: {[f'{x:.4f}' for x in embedding[:5]]}")
        print(f"Last 5 values: {[f'{x:.4f}' for x in embedding[-5:]]}")
        print("✓ Success")
        return True
    except Exception as e:
        print(f"✗ Error: {e}")
        return False


def test_error_handling():
    """Test error handling."""
    print("\n" + "=" * 60)
    print("Testing Error Handling")
    print("=" * 60)

    # Test with invalid project ID
    print("\nTest 1: Invalid project ID")
    try:
        client = VertexAIEmbedding(
            project_id="invalid-project-id-12345", location="us-central1"
        )
        embedding = client.get_text_embedding("test")
        print("✗ Expected error but got success")
        return False
    except RuntimeError as e:
        print(f"✓ Correctly caught error: {str(e)[:100]}...")

    # Test with empty text
    print("\nTest 2: Empty text")
    project_id = os.environ.get("GCP_PROJECT_ID", "your-project-id")
    try:
        client = VertexAIEmbedding(project_id=project_id)
        embedding = client.get_text_embedding("")
        print(f"⚠ Got embedding for empty text (dimension: {len(embedding)})")
    except RuntimeError as e:
        print(f"✓ Correctly caught error: {str(e)[:100]}...")

    return True


def main():
    """Run all tests."""
    print("Vertex AI Embedding Library Test Suite")
    print("=" * 60)
    print()

    # while True:
    #     import time
    #     time.sleep(1000)

    # Check if required environment variables are set
    if os.environ.get("GCP_PROJECT_ID") == "your-project-id" or not os.environ.get(
        "GCP_PROJECT_ID"
    ):
        print("⚠ WARNING: Please set the following environment variables:")
        print("  - GCP_PROJECT_ID: Your GCP project ID")
        print("  - GOOGLE_APPLICATION_CREDENTIALS: Path to credentials JSON (optional)")
        print("  - TEST_IMAGE_PATH: Path to test image (optional)")
        print("\nExample:")
        print("  export GCP_PROJECT_ID='my-project-id'")
        print("  export GOOGLE_APPLICATION_CREDENTIALS='/path/to/credentials.json'")
        print("  export TEST_IMAGE_PATH='/path/to/image.jpg'")
        print("\nContinuing with default values (tests may fail)...")
        print()

    results = []

    # Run tests
    results.append(("Text Embedding", test_text_embedding()))
    results.append(("Image Embedding", test_image_embedding()))
    results.append(("Error Handling", test_error_handling()))

    # Print summary
    print("\n" + "=" * 60)
    print("Test Summary")
    print("=" * 60)

    for test_name, passed in results:
        status = "✓ PASSED" if passed else "✗ FAILED"
        print(f"{test_name:.<50} {status}")

    print()

    # Exit with appropriate code
    if all(result for _, result in results):
        print("All tests passed!")
        sys.exit(0)
    else:
        print("Some tests failed.")
        sys.exit(1)


if __name__ == "__main__":
    main()
