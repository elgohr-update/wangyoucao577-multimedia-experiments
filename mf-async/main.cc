
#include "async_callback.h"
#include "square_root.h"
#include <iostream>

class MySqrRootCallback : public AsyncCallback {
  HANDLE event_;
  SquareRoot *sqr_root_;
  double val_;

  HRESULT hr_status_;

public:
  MySqrRootCallback(SquareRoot *sqr_root, HRESULT *hr)
      : sqr_root_(sqr_root), hr_status_(E_PENDING) {
    *hr = S_OK;

    sqr_root_->AddRef();

    event_ = CreateEvent(NULL, FALSE, FALSE, NULL);

    if (event_ == NULL) {
      *hr = HRESULT_FROM_WIN32(GetLastError());
    }
  }
  ~MySqrRootCallback() {
    sqr_root_->Release();
    CloseHandle(event_);
  }

  HRESULT WaitForCompletion(DWORD msec) {
    DWORD result = WaitForSingleObject(event_, msec);

    switch (result) {
    case WAIT_TIMEOUT:
      return E_PENDING;

    case WAIT_ABANDONED:
    case WAIT_OBJECT_0:
      return hr_status_;

    default:
      return HRESULT_FROM_WIN32(GetLastError());
    }
  }

  double value() const { return val_; }

  STDMETHODIMP Invoke(IMFAsyncResult *pResult) override {
    hr_status_ = sqr_root_->EndSquareRoot(pResult, &val_);

    SetEvent(event_);

    return S_OK;
  }
};

int main() {
  auto hr = MFStartup(MF_VERSION);
  if (FAILED(hr)) {
    std::cout << "MFStartup err " << hr << std::endl;
    return -1;
  }

  double x = 1.1;
  SquareRoot s;

  auto pCB = new (std::nothrow) MySqrRootCallback(&s, &hr);
  if (pCB == NULL) {
    hr = E_OUTOFMEMORY;
  }
  if (FAILED(hr)) {
    std::cout << "create MySqrRootCallback failed, err " << hr << std::endl;
    return -1;
  }

  // Start an asynchronous request to read data.
  hr = s.BeginSquareRoot(x, pCB, nullptr);
  if (FAILED(hr)) {
    std::cout << "BeginSquareRoot err " << hr << std::endl;
    return -1;
  }

  // wait for completion if need
  hr = pCB->WaitForCompletion(1000000);
  if (FAILED(hr)) {
    std::cout << "WaitForCompletion err " << hr << std::endl;
    return -1;
  }

  std::cout << "square_root(" << x << ") = " << pCB->value() << std::endl;

  MFShutdown();
  return 0;
}