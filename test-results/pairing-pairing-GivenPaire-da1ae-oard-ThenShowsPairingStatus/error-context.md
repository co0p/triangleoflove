# Page snapshot

```yaml
- generic [ref=e2]:
  - navigation [ref=e3]:
    - link "Go to dashboard" [ref=e5] [cursor=pointer]:
      - /url: /dashboard
      - img "Triangle of Love" [ref=e6]
    - link "Go to profile" [ref=e7] [cursor=pointer]:
      - /url: /profile
      - generic [ref=e8]: Hello, River
      - generic [ref=e9]: R
  - generic [ref=e10]:
    - banner [ref=e11]:
      - heading "Welcome back, River" [level=1] [ref=e12]
    - main [ref=e13]:
      - generic [ref=e14]:
        - paragraph [ref=e17]: Not checked in yet
        - link "Daily check-in" [ref=e18] [cursor=pointer]:
          - /url: /checkin
      - generic [ref=e19]:
        - paragraph [ref=e22]: Not connected yet
        - link "Connect with partner" [ref=e23] [cursor=pointer]:
          - /url: /pairing
```