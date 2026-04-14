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
  - main [ref=e11]:
    - heading "Pairing" [level=1] [ref=e12]
    - generic [ref=e13]:
      - generic [ref=e14]:
        - paragraph [ref=e15]: Your invite code
        - paragraph [ref=e16]: UKLH1X
      - button "Regenerate" [ref=e17] [cursor=pointer]
    - generic [ref=e18]:
      - generic [ref=e19]:
        - generic [ref=e20]: Partner's code
        - textbox "Partner's code" [ref=e21]:
          - /placeholder: Enter 6-character code
      - button "Connect" [ref=e22] [cursor=pointer]
```